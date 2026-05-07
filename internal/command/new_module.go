package command

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/charmbracelet/huh"
	"github.com/linkeunid/ligo-cli/internal/schematic/ext"
	"github.com/linkeunid/ligo-cli/internal/templateutil"
	"github.com/spf13/cobra"
)

var (
	newModuleModuleFlag string
	newModuleNoGitFlag  bool
)

var newModuleCmd = &cobra.Command{
	Use:     "new-module <name>",
	Aliases: []string{"nm"},
	Short:   "Scaffold a new ligo extension module",
	Args:    cobra.MaximumNArgs(1),
	RunE:    runNewModule,
}

func init() {
	newModuleCmd.Flags().StringVar(&newModuleModuleFlag, "module", "", "Go module path (default: github.com/<name>)")
	newModuleCmd.Flags().BoolVar(&newModuleNoGitFlag, "no-git", false, "Skip git init")
	rootCmd.AddCommand(newModuleCmd)
}

func runNewModule(cmd *cobra.Command, args []string) error {
	var name string
	var modulePath string

	if len(args) == 0 {
		// Interactive mode
		if err := promptForModule(&name, &modulePath); err != nil {
			return err
		}
	} else {
		name = args[0]
		modulePath = newModuleModuleFlag
		if modulePath == "" {
			modulePath = "github.com/" + name
		}
	}

	return scaffoldExtensionModule(name, modulePath)
}

func promptForModule(name, modulePath *string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Extension name").
				Value(name).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Go module path").
				Placeholder("github.com/myext").
				Value(modulePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("module path cannot be empty")
					}
					return nil
				}),
		),
	)
	return form.Run()
}

func scaffoldExtensionModule(name, modulePath string) error {
	names := templateutil.NormalizeName(name)
	data := map[string]string{
		"ModulePath": modulePath,
		"PascalName": names.Pascal,
		"SnakeName":  names.Snake,
	}

	if err := os.MkdirAll(name, 0755); err != nil {
		return err
	}

	err := fs.WalkDir(ext.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == "." || path == "fs.go" {
			return nil
		}

		rel := strings.TrimSuffix(path, ".tmpl")
		// Convert readme.md to README.md
		if rel == "readme.md" {
			rel = "README.md"
		}
		dest := filepath.Join(name, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		content, err := ext.FS.ReadFile(path)
		if err != nil {
			return err
		}

		t, err := template.New(path).Option("missingkey=zero").Parse(string(content))
		if err != nil {
			return writeFile(dest, content)
		}

		f, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer f.Close()

		if execErr := t.Execute(f, data); execErr != nil {
			f.Close()
			return writeFile(dest, content)
		}

		fmt.Printf("  CREATE %s\n", dest)
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("\nRunning go mod tidy...")
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = name
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: go mod tidy failed:", err)
	}

	if !newModuleNoGitFlag {
		fmt.Println("Running git init...")
		gitInit := exec.Command("git", "init")
		gitInit.Dir = name
		gitInit.Stdout = os.Stdout
		gitInit.Stderr = os.Stderr
		_ = gitInit.Run()
	}

	fmt.Printf("\nDone! cd %s && start building your extension\n", name)
	return nil
}
