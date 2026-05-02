package command

import (
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is "dev" by default; overridden by ldflags or read from build info.
var version = "dev"

func resolveVersion() string {
	if version != "dev" {
		return version // ldflags took precedence
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

var rootCmd = &cobra.Command{
	Use:   "ligo",
	Short: "Ligo framework CLI",
	Long:  "A CLI for scaffolding and managing Ligo framework projects.",
}

func Execute() {
	rootCmd.Version = resolveVersion()
	// Add -v shorthand to Cobra's auto-generated --version flag.
	rootCmd.InitDefaultVersionFlag()
	if f := rootCmd.Flags().Lookup("version"); f != nil {
		f.Shorthand = "v"
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
