package command

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ligo",
	Short: "Ligo framework CLI",
	Long:  "A CLI for scaffolding and managing Ligo framework projects.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
