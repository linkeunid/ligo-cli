// Package wired generates wire-style static DI initializers for Ligo
// applications.
//
// The user writes an "injector" file guarded by `//go:build wireinject`
// containing a function whose body calls a local `wire(...)` marker:
//
//	//go:build wireinject
//
//	package wired
//
//	import "github.com/example/app/internal/wired/factories"
//
//	func Build() (*factories.App, error) {
//	    return nil, wire(
//	        factories.NewConfig,
//	        factories.NewLogger,
//	        factories.NewDatabase,
//	        factories.NewService,
//	        factories.NewApp,
//	    )
//	}
//
//	// wire is a stub. ligo-cli replaces this whole file's Build body in
//	// the generated counterpart and never invokes wire at runtime.
//	func wire(_ ...any) error { return nil }
//
// After `ligo g wired` runs, a sibling `wired_gen.go` (built without the
// `wireinject` tag) is emitted that replaces Build's body with a hand-wired
// dependency chain in topological order. The last factory in the wire(...)
// call list is the "root" — its return value becomes the function's first
// return value, paired with a nil error.
//
// Design constraints:
//   - No hardcoded ligo import path. Factories may return any type, and the
//     generator never references "github.com/linkeunid/ligo".
//   - The marker function name is fixed (`wire`) so the parser stays simple.
//     Users define it as a stub returning `error` in the wireinject file.
//   - Only single-return or (value, error) factories are supported.
package wired

import "go/token"

// Spec describes one injector function and the providers it composes.
type Spec struct {
	PkgPath   string
	PkgName   string
	OutFile   string
	Injector  InjectorSpec
	Providers []ProviderSpec
	Imports   []ImportSpec
}

// InjectorSpec is the function the generator replaces.
type InjectorSpec struct {
	Name    string
	Params  []ParamSpec
	Results []ParamSpec
	Pos     token.Position
}

// ParamSpec is one parameter or result of the injector signature.
type ParamSpec struct {
	Name string
	Type string
}

// ProviderSpec is one factory call extracted from wire(...).
type ProviderSpec struct {
	VarName      string
	FuncRef      string
	Params       []string
	ReturnType   string
	ReturnsError bool
	DepVars      []string
}

// ImportSpec captures a package the emitted file needs.
type ImportSpec struct {
	Alias string
	Path  string
}
