package schematic

import "fmt"

// Context holds all data needed by a schematic to generate files.
type Context struct {
	Name       string // normalized input (e.g. "product")
	ModulePath string // from go.mod (e.g. "github.com/acme/my-app")
	WorkDir    string // project root (cwd)
	WithAuth   bool   // --with-auth flag
	Full       bool   // --full flag (use full templates with repo/DTOs)
	DryRun     bool   // --dry-run flag (print actions without writing files)
}

// Schematic generates one or more files into a ligo project.
type Schematic interface {
	Name() string
	Aliases() []string
	Run(ctx Context) error
}

var registry = map[string]Schematic{}

// Register adds a schematic (and its aliases) to the global registry.
func Register(s Schematic) {
	registry[s.Name()] = s
	for _, a := range s.Aliases() {
		registry[a] = s
	}
}

// Lookup returns the schematic for the given name or alias.
func Lookup(name string) (Schematic, error) {
	s, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown schematic %q — available: resource, module (mo), controller (co), usecase (uc), entity, repository, dto, presenter", name)
	}
	return s, nil
}
