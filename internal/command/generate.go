package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/linkeunid/ligo-cli/internal/project"
	"github.com/linkeunid/ligo-cli/internal/schematic"
)

var (
	withAuth      bool
	dryRun        bool
	fullSchematic bool
)

var generateCmd = &cobra.Command{
	Use:     "generate <schematic> [name]",
	Aliases: []string{"g"},
	Short:   "Generate a schematic into the current ligo project",
	Args:    cobra.RangeArgs(1, 2),
	RunE:    runGenerate,
}

func init() {
	generateCmd.Flags().BoolVar(&withAuth, "with-auth", false, "Wire AuthGuard into generated controller (requires --full)")
	generateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview files without writing them")
	generateCmd.Flags().BoolVar(&fullSchematic, "full", false, "Generate full template with repo/DTOs")
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	schematicName := args[0]
	name := ""
	if len(args) == 2 {
		name = args[1]
	}

	info, err := project.Detect(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	s, err := schematic.Lookup(schematicName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	return s.Run(schematic.Context{
		Name:       name,
		ModulePath: info.ModulePath,
		WorkDir:    ".",
		WithAuth:   withAuth,
		Full:       fullSchematic,
		DryRun:     dryRun,
	})
}
