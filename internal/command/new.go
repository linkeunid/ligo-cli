package command

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/linkeunid/ligo-cli/internal/newtemplate"
	"github.com/spf13/cobra"
)

var (
	newModuleFlag string
	noGitFlag     bool
)

var newCmd = &cobra.Command{
	Use:     "new <project-name>",
	Aliases: []string{"n"},
	Short:   "Scaffold a new ligo project",
	Args:    cobra.ExactArgs(1),
	RunE:    runNew,
}

func init() {
	newCmd.Flags().StringVar(&newModuleFlag, "module", "", "Go module path (default: github.com/<project-name>)")
	newCmd.Flags().BoolVar(&noGitFlag, "no-git", false, "Skip git init")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	modulePath := newModuleFlag
	if modulePath == "" {
		modulePath = "github.com/" + projectName
	}

	data := map[string]string{"ModulePath": modulePath}

	if err := os.MkdirAll(projectName, 0755); err != nil {
		return err
	}

	err := fs.WalkDir(newtemplate.FS, "testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(path, "testdata/")
		if rel == "" || rel == "testdata" {
			return nil
		}
		// Restore .tmpl suffix stripped from go.mod and go.sum to avoid module boundary issues
		rel = strings.TrimSuffix(rel, ".tmpl")
		dest := filepath.Join(projectName, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		content, err := newtemplate.FS.ReadFile(path)
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
	tidy.Dir = projectName
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	if err := tidy.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: go mod tidy failed:", err)
	}

	if !noGitFlag {
		fmt.Println("Running git init...")
		gitInit := exec.Command("git", "init")
		gitInit.Dir = projectName
		gitInit.Stdout = os.Stdout
		gitInit.Stderr = os.Stderr
		_ = gitInit.Run()
	}

	fmt.Printf("\nDone! cd %s && ligo serve\n", projectName)
	return nil
}

func writeFile(dest string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, content, 0644); err != nil {
		return err
	}
	fmt.Printf("  CREATE %s\n", dest)
	return nil
}
