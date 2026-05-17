package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type presenterSchematic struct{}

func init() { Register(&presenterSchematic{}) }

func (s *presenterSchematic) Name() string      { return "presenter" }
func (s *presenterSchematic) Aliases() []string { return []string{"pre"} }

func (s *presenterSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	return templateutil.RenderToFile(
		presenterTmpl,
		filepath.Join(ctx.WorkDir, "internal", "infrastructure", "http", "presenter",
			fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	)
}
