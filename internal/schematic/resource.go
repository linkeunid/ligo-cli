package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type resourceSchematic struct{}

func init() { Register(&resourceSchematic{}) }

func (s *resourceSchematic) Name() string      { return "resource" }
func (s *resourceSchematic) Aliases() []string { return []string{"res"} }

func (s *resourceSchematic) Run(ctx Context) error {
	if ctx.Name == "" {
		var name string
		var withAuth bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Resource name:").Value(&name),
				huh.NewConfirm().Title("Generate with auth guard?").Value(&withAuth),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
		ctx.Name = name
		ctx.WithAuth = withAuth
	}

	n := templateutil.NormalizeName(ctx.Name)
	d := makeData(ctx)

	steps := []struct {
		label string
		fn    func() error
	}{
		{"entity", func() error { return (&entitySchematic{}).Run(ctx) }},
		{"dto", func() error { return (&dtoSchematic{}).Run(ctx) }},
		{"usecase", func() error {
			return templateutil.RenderToFile(usecaseTmpl,
				filepath.Join(ctx.WorkDir, "internal", "usecase", fmt.Sprintf("%s.go", n.Snake)),
				d, ctx.DryRun)
		}},
		{"repository", func() error { return (&repositorySchematic{}).Run(ctx) }},
		{"controller", func() error { return (&controllerSchematic{}).Run(ctx) }},
		{"presenter", func() error { return (&presenterSchematic{}).Run(ctx) }},
		{"module", func() error {
			return templateutil.RenderToFile(moduleTmpl,
				filepath.Join(ctx.WorkDir, "internal", "module", fmt.Sprintf("%s.go", n.Snake)),
				d, ctx.DryRun)
		}},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.label, err)
		}
	}

	return askAndRegister(ctx, n.Pascal)
}
