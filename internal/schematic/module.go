package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type moduleSchematic struct{}

func init() { Register(&moduleSchematic{}) }

func (s *moduleSchematic) Name() string      { return "module" }
func (s *moduleSchematic) Aliases() []string { return []string{"mo"} }

func (s *moduleSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	tmpl := simpleModuleTmpl
	if ctx.Full {
		tmpl = moduleTmpl
	}

	if err := templateutil.RenderToFile(
		tmpl,
		filepath.Join(ctx.WorkDir, "internal", "module", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	return askAndRegister(ctx, n.Pascal)
}

// askAndRegister prompts to register the module in module/main.go.
func askAndRegister(ctx Context, pascal string) error {
	var register bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Register %s module in internal/module/main.go?", pascal)).
				Value(&register),
		),
	)
	if err := form.Run(); err != nil {
		return err
	}
	if !register {
		return nil
	}

	mainPath := filepath.Join(ctx.WorkDir, "internal", "module", "main.go")
	insertion := fmt.Sprintf("\t\t\t%sModule(),\n", pascal)
	return templateutil.PatchFile(mainPath, "\t\t),\n", insertion, ctx.DryRun)
}
