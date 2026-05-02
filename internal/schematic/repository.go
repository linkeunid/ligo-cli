package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type repositorySchematic struct{}

func init() { Register(&repositorySchematic{}) }

func (s *repositorySchematic) Name() string      { return "repository" }
func (s *repositorySchematic) Aliases() []string { return []string{"rep"} }

func (s *repositorySchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	return templateutil.RenderToFile(memoryRepoTmpl,
		filepath.Join(ctx.WorkDir, "internal", "infrastructure", "persistence", "memory",
			fmt.Sprintf("%s_repo.go", n.Snake)),
		d, ctx.DryRun,
	)
}
