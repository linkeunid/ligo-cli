package wired

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// LoadOptions controls how Load resolves the target package.
type LoadOptions struct {
	// PkgPattern is a go-list pattern selecting exactly one package.
	PkgPattern string
	// Dir is the working directory for package loading. Defaults to cwd.
	Dir string
	// InjectorName names the function whose body to rewrite. Default "Build".
	InjectorName string
	// MarkerName names the user-defined stub function whose arguments
	// enumerate the factories. Default "wire".
	MarkerName string
	// OutFile is the destination filename written alongside the package.
	// Default "wired_gen.go".
	OutFile string
}

// Load parses the target package under the wireinject build tag, locates the
// injector and marker calls, and resolves every factory's signature into a
// topologically sorted Spec ready for Emit.
func Load(opts LoadOptions) (*Spec, error) {
	if opts.PkgPattern == "" {
		return nil, fmt.Errorf("wired: PkgPattern is required")
	}
	if opts.InjectorName == "" {
		opts.InjectorName = "Build"
	}
	if opts.MarkerName == "" {
		opts.MarkerName = "wire"
	}
	if opts.OutFile == "" {
		opts.OutFile = "wired_gen.go"
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax |
			packages.NeedTypesInfo | packages.NeedDeps,
		Dir:        opts.Dir,
		BuildFlags: []string{"-tags=wireinject"},
	}
	pkgs, err := packages.Load(cfg, opts.PkgPattern)
	if err != nil {
		return nil, fmt.Errorf("wired: load package: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("wired: pattern %q matched no packages", opts.PkgPattern)
	}
	if len(pkgs) > 1 {
		return nil, fmt.Errorf("wired: pattern %q matched %d packages, expected 1",
			opts.PkgPattern, len(pkgs))
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("wired: package %q errors: %v", pkg.PkgPath, pkg.Errors[0])
	}

	fn, file, err := findInjector(pkg, opts.InjectorName)
	if err != nil {
		return nil, err
	}

	factoryArgs, err := extractMarkerArgs(fn, opts.MarkerName)
	if err != nil {
		return nil, err
	}

	imports := map[string]string{}
	providers, err := buildProviders(pkg, factoryArgs, imports)
	if err != nil {
		return nil, err
	}

	injector, err := buildInjector(pkg, fn, imports)
	if err != nil {
		return nil, err
	}

	ordered, err := topoSort(providers, injector.Params)
	if err != nil {
		return nil, err
	}

	spec := &Spec{
		PkgPath:   pkg.PkgPath,
		PkgName:   pkg.Name,
		OutFile:   filepath.Join(packageDir(pkg), opts.OutFile),
		Injector:  injector,
		Providers: ordered,
		Imports:   sortedImports(imports, pkg.PkgPath),
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	_ = file
	return spec, nil
}

// Validate sanity-checks the spec before emission.
func (s *Spec) Validate() error {
	if s.PkgName == "" {
		return fmt.Errorf("wired: spec missing package name")
	}
	if s.Injector.Name == "" {
		return fmt.Errorf("wired: spec missing injector name")
	}
	if len(s.Providers) == 0 {
		return fmt.Errorf("wired: no providers extracted from %s", s.Injector.Name)
	}
	for i, p := range s.Providers {
		if len(p.DepVars) != len(p.Params) {
			return fmt.Errorf("wired: provider #%d (%s) has %d params but %d resolved deps",
				i, p.FuncRef, len(p.Params), len(p.DepVars))
		}
	}
	return nil
}

func packageDir(pkg *packages.Package) string {
	files := pkg.GoFiles
	if len(files) == 0 {
		files = pkg.CompiledGoFiles
	}
	if len(files) == 0 {
		return "."
	}
	return filepath.Dir(files[0])
}

func findInjector(pkg *packages.Package, name string) (*ast.FuncDecl, *ast.File, error) {
	for _, file := range pkg.Syntax {
		if !hasWireinjectTag(file) {
			continue
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != name {
				continue
			}
			return fn, file, nil
		}
	}
	return nil, nil, fmt.Errorf("wired: injector %q not found in any //go:build wireinject file of %s",
		name, pkg.PkgPath)
}

func hasWireinjectTag(file *ast.File) bool {
	for _, group := range file.Comments {
		if group.Pos() >= file.Package {
			continue
		}
		for _, c := range group.List {
			if strings.HasPrefix(c.Text, "//go:build") && strings.Contains(c.Text, "wireinject") {
				return true
			}
		}
	}
	return false
}

func extractMarkerArgs(fn *ast.FuncDecl, marker string) ([]ast.Expr, error) {
	if fn.Body == nil {
		return nil, fmt.Errorf("wired: injector %q has no body", fn.Name.Name)
	}
	var args []ast.Expr
	var found bool
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok || ident.Name != marker {
			return true
		}
		args = call.Args
		found = true
		return false
	})
	if !found {
		return nil, fmt.Errorf("wired: injector %q body must call %s(factory1, factory2, ...)",
			fn.Name.Name, marker)
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("wired: %s() needs at least one factory argument", marker)
	}
	return args, nil
}

func buildProviders(pkg *packages.Package, args []ast.Expr, imports map[string]string) ([]ProviderSpec, error) {
	out := make([]ProviderSpec, 0, len(args))
	for _, arg := range args {
		tv, ok := pkg.TypesInfo.Types[arg]
		if !ok {
			return nil, fmt.Errorf("wired: cannot type-check factory %s", exprText(arg))
		}
		sig, ok := tv.Type.(*types.Signature)
		if !ok {
			return nil, fmt.Errorf("wired: argument %s is not a function value", exprText(arg))
		}
		if sig.Results().Len() == 0 {
			return nil, fmt.Errorf("wired: factory %s must return at least one value", exprText(arg))
		}
		if sig.Results().Len() > 2 {
			return nil, fmt.Errorf("wired: factory %s returns %d values; want 1 or (value, error)",
				exprText(arg), sig.Results().Len())
		}
		primary := sig.Results().At(0).Type()
		returnsError := false
		if sig.Results().Len() == 2 {
			second := sig.Results().At(1).Type()
			if !isErrorType(second) {
				return nil, fmt.Errorf("wired: factory %s second result must be error, got %s",
					exprText(arg), second.String())
			}
			returnsError = true
		}

		ps := ProviderSpec{
			FuncRef:      renderFuncRef(pkg, arg, imports),
			ReturnType:   qualifyType(pkg, primary, imports),
			ReturnsError: returnsError,
			VarName:      varNameFromType(primary),
		}
		for i := 0; i < sig.Params().Len(); i++ {
			ps.Params = append(ps.Params, qualifyType(pkg, sig.Params().At(i).Type(), imports))
		}
		out = append(out, ps)
	}
	return out, nil
}

func isErrorType(t types.Type) bool {
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return named.Obj().Pkg() == nil && named.Obj().Name() == "error"
}

func renderFuncRef(pkg *packages.Package, expr ast.Expr, imports map[string]string) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			if obj := pkg.TypesInfo.ObjectOf(ident); obj != nil {
				if pn, ok := obj.(*types.PkgName); ok {
					recordImport(imports, pn.Imported().Path(), pn.Name())
					return fmt.Sprintf("%s.%s", pn.Name(), e.Sel.Name)
				}
			}
		}
		return exprText(e)
	default:
		return exprText(expr)
	}
}

func recordImport(imports map[string]string, path, alias string) {
	if path == "" {
		return
	}
	if existing, ok := imports[path]; ok && existing != "" {
		return
	}
	imports[path] = alias
}

func qualifyType(pkg *packages.Package, t types.Type, imports map[string]string) string {
	qualifier := func(p *types.Package) string {
		if p == nil || p.Path() == pkg.PkgPath {
			return ""
		}
		recordImport(imports, p.Path(), p.Name())
		return p.Name()
	}
	return types.TypeString(t, qualifier)
}

func exprText(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprText(e.X) + "." + e.Sel.Name
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func buildInjector(pkg *packages.Package, fn *ast.FuncDecl, imports map[string]string) (InjectorSpec, error) {
	spec := InjectorSpec{
		Name: fn.Name.Name,
		Pos:  pkg.Fset.Position(fn.Pos()),
	}
	if fn.Type.Params != nil {
		for _, field := range fn.Type.Params.List {
			typeStr := qualifyType(pkg, pkg.TypesInfo.TypeOf(field.Type), imports)
			names := fieldNames(field)
			for _, n := range names {
				spec.Params = append(spec.Params, ParamSpec{Name: n, Type: typeStr})
			}
		}
	}
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			typeStr := qualifyType(pkg, pkg.TypesInfo.TypeOf(field.Type), imports)
			names := fieldNames(field)
			if len(names) == 0 {
				spec.Results = append(spec.Results, ParamSpec{Type: typeStr})
				continue
			}
			for _, n := range names {
				spec.Results = append(spec.Results, ParamSpec{Name: n, Type: typeStr})
			}
		}
	}
	return spec, nil
}

func fieldNames(field *ast.Field) []string {
	if len(field.Names) == 0 {
		return nil
	}
	out := make([]string, len(field.Names))
	for i, n := range field.Names {
		out[i] = n.Name
	}
	return out
}

func varNameFromType(t types.Type) string {
	stripped := t
	if ptr, ok := stripped.(*types.Pointer); ok {
		stripped = ptr.Elem()
	}
	named, ok := stripped.(*types.Named)
	if !ok {
		return sanitizeIdent("v")
	}
	name := named.Obj().Name()
	if name == "" {
		return sanitizeIdent("v")
	}
	first := rune(name[0])
	return sanitizeIdent(string(unicode.ToLower(first)) + name[1:])
}

func sanitizeIdent(s string) string {
	if s == "" {
		return "v"
	}
	switch s {
	case "type", "func", "var", "const", "import", "package", "interface", "map", "chan", "go", "return", "if", "for", "select", "switch", "case", "default", "range", "break", "continue", "fallthrough", "goto", "defer", "error":
		return s + "_"
	}
	return s
}

func sortedImports(m map[string]string, selfPath string) []ImportSpec {
	out := make([]ImportSpec, 0, len(m))
	for path, alias := range m {
		if path == selfPath {
			continue
		}
		out = append(out, ImportSpec{Alias: alias, Path: path})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
