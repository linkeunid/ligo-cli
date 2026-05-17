package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/linkeunid/ligo-cli/internal/project"
)

func writeGoMod(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetect_ValidLigoProject(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, `module github.com/acme/my-app

go 1.22

require github.com/linkeunid/ligo v0.1.0
`)
	info, err := project.Detect(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ModulePath != "github.com/acme/my-app" {
		t.Errorf("ModulePath got %q want %q", info.ModulePath, "github.com/acme/my-app")
	}
}

func TestDetect_NotLigoProject(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, `module github.com/acme/other

go 1.22
`)
	_, err := project.Detect(dir)
	if err == nil {
		t.Fatal("expected error for non-ligo project")
	}
}

func TestDetect_NoGoMod(t *testing.T) {
	dir := t.TempDir()
	_, err := project.Detect(dir)
	if err == nil {
		t.Fatal("expected error when go.mod missing")
	}
}
