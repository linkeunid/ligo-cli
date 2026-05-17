package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var buildOut string

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build the ligo application (go build ./cmd/api/)",
	RunE:    runBuild,
}

func init() {
	buildCmd.Flags().StringVar(&buildOut, "out", "bin/app", "Output binary path")
	buildCmd.Flags().BoolVarP(&noWiredFlag, "no-wired", "n", false, "Skip auto-regeneration of wired_gen.go")
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	if err := regenerateWired(); err != nil {
		return err
	}
	c := exec.Command("go", "build", "-o", buildOut, "./cmd/api/")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Printf("Built: %s\n", buildOut)
	return nil
}
