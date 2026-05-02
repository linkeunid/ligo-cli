package schematic

import (
	_ "embed"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

//go:embed templates/generate/entity.go.tmpl
var entityTmpl string

//go:embed templates/generate/repository.go.tmpl
var repositoryTmpl string

//go:embed templates/generate/dto_create.go.tmpl
var dtoCreateTmpl string

//go:embed templates/generate/dto_update.go.tmpl
var dtoUpdateTmpl string

//go:embed templates/generate/usecase.go.tmpl
var usecaseTmpl string

//go:embed templates/generate/controller.go.tmpl
var controllerTmpl string

//go:embed templates/generate/controller_auth.go.tmpl
var controllerAuthTmpl string

//go:embed templates/generate/presenter.go.tmpl
var presenterTmpl string

//go:embed templates/generate/memory_repo.go.tmpl
var memoryRepoTmpl string

//go:embed templates/generate/memory_uuid.go.tmpl
var memoryUUIDTmpl string

//go:embed templates/generate/usecase_errors.go.tmpl
var usecaseErrorsTmpl string

//go:embed templates/generate/module.go.tmpl
var moduleTmpl string

//go:embed templates/generate/simple_usecase.go.tmpl
var simpleUsecaseTmpl string

//go:embed templates/generate/simple_controller.go.tmpl
var simpleControllerTmpl string

//go:embed templates/generate/simple_module.go.tmpl
var simpleModuleTmpl string

// TemplateData is the data passed to every generate template.
type TemplateData struct {
	ModulePath  string
	Pascal      string
	Kebab       string
	PluralKebab string
}

func makeData(ctx Context) TemplateData {
	n := templateutil.NormalizeName(ctx.Name)
	return TemplateData{
		ModulePath:  ctx.ModulePath,
		Pascal:      n.Pascal,
		Kebab:       n.Kebab,
		PluralKebab: n.PluralKebab,
	}
}
