package templateutil

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// RenderToFile renders a text/template string with data into destPath.
// Creates parent directories as needed. Prints CREATE <path> to stdout.
// If dryRun is true, prints the action without writing any files.
func RenderToFile(tmplStr, destPath string, data any, dryRun bool) error {
	if dryRun {
		fmt.Printf("  CREATE %s\n", destPath)
		return nil
	}

	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if mkErr := os.MkdirAll(filepath.Dir(destPath), 0o755); mkErr != nil {
		return fmt.Errorf("create dirs: %w", mkErr)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("render template: %w", err)
	}

	fmt.Printf("  CREATE %s\n", destPath)
	return nil
}

// PatchFile inserts a string into an existing file before the first occurrence of marker.
// Prints UPDATE <path> to stdout.
// If dryRun is true, prints the action without modifying the file.
func PatchFile(path, marker, insertion string, dryRun bool) error {
	if dryRun {
		fmt.Printf("  UPDATE %s\n", path)
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	content := string(data)
	idx := indexOf(content, marker)
	if idx == -1 {
		return fmt.Errorf("marker %q not found in %s", marker, path)
	}

	patched := content[:idx] + insertion + content[idx:]
	if err := os.WriteFile(path, []byte(patched), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	fmt.Printf("  UPDATE %s\n", path)
	return nil
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
