package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ligoImport = "github.com/linkeunid/ligo"

// Info holds information extracted from a ligo project's go.mod.
type Info struct {
	ModulePath string
	GoModPath  string
}

// Detect reads go.mod in dir and validates it is a ligo project.
func Detect(dir string) (*Info, error) {
	goModPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("go.mod not found in %s: run this command from a ligo project root", dir)
	}

	content := string(data)
	modulePath := extractModulePath(content)
	if modulePath == "" {
		return nil, fmt.Errorf("could not parse module path from go.mod")
	}

	if !strings.Contains(content, ligoImport) {
		return nil, fmt.Errorf("not a ligo project (%s not found in go.mod)", ligoImport)
	}

	return &Info{ModulePath: modulePath, GoModPath: goModPath}, nil
}

func extractModulePath(content string) string {
	for line := range strings.SplitSeq(content, "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}
