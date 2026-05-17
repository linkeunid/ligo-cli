package wired_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/linkeunid/ligo-cli/internal/schematic/wired"
)

func TestLoadAndEmit_Simple(t *testing.T) {
	spec, err := wired.Load(wired.LoadOptions{
		PkgPattern: "./",
		Dir:        "testdata/simple",
	})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got, want := len(spec.Providers), 5; got != want {
		t.Fatalf("provider count: got %d, want %d", got, want)
	}

	pos := map[string]int{}
	for i, p := range spec.Providers {
		pos[p.ReturnType] = i
	}
	mustBefore := func(a, b string) {
		t.Helper()
		ai, aok := pos[a]
		bi, bok := pos[b]
		if !aok || !bok {
			t.Fatalf("missing types in order map: %s=%v %s=%v (have %v)", a, aok, b, bok, pos)
		}
		if !(ai < bi) {
			t.Errorf("expected %s (%d) before %s (%d)", a, ai, b, bi)
		}
	}
	mustBefore("*Config", "*Database")
	mustBefore("*Logger", "*Database")
	mustBefore("*Database", "*Service")
	mustBefore("*Service", "*App")

	var buf bytes.Buffer
	if err := wired.Emit(&buf, spec); err != nil {
		t.Fatalf("Emit: %v", err)
	}

	got := buf.String()
	for _, want := range []string{
		"//go:build !wireinject",
		"package simple",
		"config := NewConfig()",
		"logger := NewLogger()",
		"database := NewDatabase(config, logger)",
		"service, err := NewService(database, logger)",
		"app := NewApp(service)",
		"return app, nil",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("emitted source missing %q\n--- got ---\n%s", want, got)
		}
	}

	// The emitted file must not import the ligo core package — the codegen
	// stays agnostic of the ligo import path.
	if strings.Contains(got, "linkeunid/ligo") {
		t.Errorf("emitted source must not reference linkeunid/ligo\n--- got ---\n%s", got)
	}
}
