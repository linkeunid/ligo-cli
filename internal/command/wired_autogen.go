package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/schematic/wired"
)

const defaultWiredPkg = "internal/wired"

// regenerateWired runs the wired codegen against the project's wired package
// when one exists. Silent no-op otherwise so projects without wired pay no
// cost. Errors are surfaced (a stale wired_gen.go will fail the build anyway,
// so it's better to know immediately than to silently launch with the old
// graph). Operates relative to the current working directory.
func regenerateWired() error {
	if noWiredFlag {
		return nil
	}
	pkgRel := defaultWiredPkg
	pkgDir := pkgRel
	if info, err := os.Stat(pkgDir); err != nil || !info.IsDir() {
		return nil
	}
	if !hasWireinjectFile(pkgDir) {
		return nil
	}

	spec, err := wired.Load(wired.LoadOptions{
		PkgPattern: "./" + pkgRel,
	})
	if err != nil {
		return fmt.Errorf("wired regen: %w", err)
	}
	outPath := filepath.Join(pkgDir, "wired_gen.go")
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("wired regen: create %s: %w", outPath, err)
	}
	defer f.Close()
	if err := wired.Emit(f, spec); err != nil {
		return fmt.Errorf("wired regen: emit: %w", err)
	}
	fmt.Printf("Regenerated %s (%d providers)\n", outPath, len(spec.Providers))
	return nil
}

func hasWireinjectFile(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".go" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if containsWireinjectBuildTag(data) {
			return true
		}
	}
	return false
}

func containsWireinjectBuildTag(src []byte) bool {
	// Build constraints must precede the package clause. Scan only the
	// header — the first "package" keyword ends the relevant window.
	const needle = "//go:build"
	for i := 0; i+len(needle) < len(src); i++ {
		if src[i] == 'p' && i+7 < len(src) && string(src[i:i+8]) == "package " {
			return false
		}
		if src[i] == '/' && i+len(needle) <= len(src) && string(src[i:i+len(needle)]) == needle {
			// Consume the rest of the line.
			end := i + len(needle)
			for end < len(src) && src[end] != '\n' {
				end++
			}
			if containsToken(string(src[i:end]), "wireinject") {
				return true
			}
			i = end
		}
	}
	return false
}

func containsToken(line, token string) bool {
	// Cheap whole-word check — avoids matching "wireinjects".
	for i := 0; i+len(token) <= len(line); i++ {
		if line[i:i+len(token)] != token {
			continue
		}
		left := i == 0 || !isWordChar(line[i-1])
		right := i+len(token) == len(line) || !isWordChar(line[i+len(token)])
		if left && right {
			return true
		}
	}
	return false
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}
