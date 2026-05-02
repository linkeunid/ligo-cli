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
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	c := exec.Command("go", "build", "-o", buildOut, "./cmd/api/")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	fmt.Printf("Built: %s\n", buildOut)
	return nil
}
