package templateutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

func TestRenderToFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out.go")

	tmpl := "package {{.Package}}\n\ntype {{.Name}} struct{}\n"
	data := map[string]string{"Package": "entity", "Name": "Product"}

	if err := templateutil.RenderToFile(tmpl, dest, data, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	want := "package entity\n\ntype Product struct{}\n"
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderToFile_CreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "a", "b", "c.go")

	if err := templateutil.RenderToFile("hello", dest, nil, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestPatchFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.go")

	original := "func Main() {\n\t\tligo.Imports(\n\t\t\tAuth(),\n\t\t),\n\t}\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := templateutil.PatchFile(path, "\t\t),\n", "\t\t\tProduct(),\n", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "func Main() {\n\t\tligo.Imports(\n\t\t\tAuth(),\n\t\t\tProduct(),\n\t\t),\n\t}\n"
	if string(got) != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
