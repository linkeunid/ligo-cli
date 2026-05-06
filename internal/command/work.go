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
	Args:    cobra.ExactArgs(1),
	RunE:    runWork,
}

func init() {
	rootCmd.AddCommand(workCmd)
}

func runWork(cmd *cobra.Command, args []string) error {
	runnerName := args[0]

	// Check if the runner exists
	runnerPath := filepath.Join("cmd", "runner", runnerName, "main.go")
	if _, err := os.Stat(runnerPath); os.IsNotExist(err) {
		return fmt.Errorf("runner %q not found at %s\nAvailable runners:", runnerName, runnerPath)
	}

	// Run the runner
	fmt.Printf("Starting runner: %s\n", runnerName)
	runCmd := exec.Command("go", "run", runnerPath)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Stdin = os.Stdin

	return runCmd.Run()
}
