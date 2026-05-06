package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var workCmd = &cobra.Command{
	Use:     "work <runner-name>",
	Short:   "Run a background worker/runner",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runWork,
}

func init() {
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

	// Run the runner
	fmt.Printf("Starting runner: %s\n", runnerName)
	runCmd := exec.Command("go", "run", runnerPath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin

	return runCmd.Run()
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
	return nil
}
