package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type controllerSchematic struct{}

func init() { Register(&controllerSchematic{}) }

func (s *controllerSchematic) Name() string      { return "controller" }
func (s *controllerSchematic) Aliases() []string { return []string{"co"} }

func (s *controllerSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	tmpl := simpleControllerTmpl
	if ctx.Full {
		tmpl = controllerTmpl
		if ctx.WithAuth {
			tmpl = controllerAuthTmpl
		}
	}

	return templateutil.RenderToFile(tmpl,
		filepath.Join(ctx.WorkDir, "internal", "infrastructure", "http", "controller",
			fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	)
}
