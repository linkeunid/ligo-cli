package schematic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/schematic/wired"
)

type wiredSchematic struct{}

func init() { Register(&wiredSchematic{}) }

func (s *wiredSchematic) Name() string      { return "wired" }
func (s *wiredSchematic) Aliases() []string { return []string{"wire"} }

// Run scans `internal/wired` (or ctx.Name as a package path relative to the
// project root) for a //go:build wireinject injector and emits a parallel
// wired_gen.go alongside it.
func (s *wiredSchematic) Run(ctx Context) error {
	pkgRel := ctx.Name
	if pkgRel == "" {
		pkgRel = "internal/wired"
	}

	pkgDir := filepath.Join(ctx.WorkDir, pkgRel)
	info, err := os.Stat(pkgDir)
	if err != nil {
		return fmt.Errorf("wired: package directory %q not found: %w", pkgDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("wired: %q is not a directory", pkgDir)
	}

	spec, err := wired.Load(wired.LoadOptions{
		PkgPattern: "./" + pkgRel,
		Dir:        ctx.WorkDir,
	})
	if err != nil {
		return err
	}

	outPath := filepath.Join(pkgDir, "wired_gen.go")

	if ctx.DryRun {
		fmt.Printf("CREATE %s (would wire %d providers, root=%s.%s)\n",
			outPath, len(spec.Providers), spec.PkgPath, spec.Injector.Name)
		return nil
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("wired: create %s: %w", outPath, err)
	}
	defer f.Close()

	if err := wired.Emit(f, spec); err != nil {
		return fmt.Errorf("wired: emit: %w", err)
	}

	fmt.Printf("CREATE %s\n", outPath)
	fmt.Printf("  injector : %s.%s\n", spec.PkgPath, spec.Injector.Name)
	fmt.Printf("  providers: %d\n", len(spec.Providers))
	return nil
}
