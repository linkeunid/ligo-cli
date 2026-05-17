package command

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var workWatchFlag bool

var workCmd = &cobra.Command{
	Use:   "work <runner-name>",
	Short: "Run a background worker/runner",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runWork,
}

func init() {
	workCmd.Flags().BoolVarP(&workWatchFlag, "watch", "w", false, "Restart on file changes")
	rootCmd.AddCommand(workCmd)
}

func runWork(cmd *cobra.Command, args []string) error {
	// If no runner name provided, list available runners
	if len(args) == 0 {
		return listRunners()
	}

	runnerName := args[0]

	// Check if the runner exists
	runnerPath := filepath.Join("cmd", "runner", runnerName, "main.go")
	if _, err := os.Stat(runnerPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: runner %q not found\n\n", runnerName)
		listRunners()
		return fmt.Errorf("runner %q not found", runnerName)
	}

	if workWatchFlag {
		return runWorkWithWatcher(runnerName, runnerPath)
	}

	// Run the runner
	fmt.Printf("Starting runner: %s\n", runnerName)
	runCmd := exec.Command("go", "run", runnerPath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin

	return runCmd.Run()
}

func runWorkWithWatcher(runnerName, runnerPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	defer watcher.Close()

	// Watch the runner's internal directory and cmd directory
	watchDirs := []string{
		filepath.Join("cmd", "runner", runnerName),
		filepath.Join("internal", runnerName),
	}
	for _, dir := range watchDirs {
		if err := addRecursive(watcher, dir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not watch %s: %v\n", dir, err)
		}
	}

	var proc *exec.Cmd
	startApp := func() {
		proc = exec.Command("go", "run", runnerPath)
		proc.Stdout = os.Stdout
		proc.Stderr = os.Stderr
		proc.Stdin = os.Stdin
		proc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := proc.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Error starting runner:", err)
		}
	}

	killApp := func() {
		if proc != nil && proc.Process != nil {
			// Kill the entire process group so the go run child (the runner) is also killed.
			_ = syscall.Kill(-proc.Process.Pid, syscall.SIGKILL)
			_ = proc.Wait()
			proc = nil
		}
	}

	fmt.Printf("Starting runner: %s (with --watch)\n", runnerName)
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
			fmt.Printf("File changed — restarting runner: %s\n", runnerName)
			killApp()
			startApp()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintln(os.Stderr, "Watcher error:", err)
		case <-quit:
			if proc != nil && proc.Process != nil {
				// Forward SIGINT to the child's process group so it shuts down gracefully.
				_ = syscall.Kill(-proc.Process.Pid, syscall.SIGINT)
				_ = proc.Wait()
			}
			return nil
		}
	}
}

func listRunners() error {
	runnerDir := filepath.Join("cmd", "runner")

	entries, err := os.ReadDir(runnerDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No runners found. Create one with:")
			fmt.Println("  ligo g run <runner-name>")
			return nil
		}
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No runners found. Create one with:")
		fmt.Println("  ligo g run <runner-name>")
		return nil
	}

	fmt.Println("Available runners:")
	for _, entry := range entries {
		if entry.IsDir() {
			mainPath := filepath.Join(runnerDir, entry.Name(), "main.go")
			if _, err := os.Stat(mainPath); err == nil {
				fmt.Printf("  - %s\n", entry.Name())
			}
		}
	}

	fmt.Println("\nRun a runner:")
	fmt.Println("  ligo work <runner-name>")
	fmt.Println("  ligo work <runner-name> --watch    # with auto-reload on file changes")
	return nil
}
