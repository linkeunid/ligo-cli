package command

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchFlag bool

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Run the ligo application (go run ./cmd/api/)",
	RunE:    runServe,
}

func init() {
	serveCmd.Flags().BoolVar(&watchFlag, "watch", false, "Restart on file changes")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	if !watchFlag {
		c := exec.Command("go", "run", "./cmd/api/")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	}
	return runWithWatcher()
}

func runWithWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	watchDirs := []string{"internal", "cmd"}
	for _, dir := range watchDirs {
		if err := addRecursive(watcher, dir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not watch %s: %v\n", dir, err)
		}
	}

	var proc *exec.Cmd
	startApp := func() {
		proc = exec.Command("go", "run", "./cmd/api/")
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		proc.Stdin = os.Stdin
		proc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := proc.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error starting app:", err)
		}
	}

	killApp := func() {
		if proc != nil && proc.Process != nil {
			// Kill the entire process group so the go run child (the app) is also killed.
			_ = syscall.Kill(-proc.Process.Pid, syscall.SIGKILL)
			_ = proc.Wait()
			proc = nil
		}
	}

	startApp()

	debounce := time.NewTimer(0)
	<-debounce.C

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if filepath.Ext(event.Name) == ".go" {
				debounce.Reset(500 * time.Millisecond)
			}
		case <-debounce.C:
			clearScreen()
			fmt.Println("File changed — restarting...")
			killApp()
			startApp()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintln(os.Stderr, "Watcher error:", err)
		case <-quit:
			killApp()
			return nil
		}
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func addRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return w.Add(path)
		}
		return nil
	})
}
