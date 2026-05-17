//go:build wireinject

package simple

func Build() (*App, error) {
	return nil, wire(
		NewConfig,
		NewLogger,
		NewDatabase,
		NewService,
		NewApp,
	)
}

// wire is the user-defined stub. ligo-cli replaces Build's body in the
// generated counterpart; this function is never called at runtime.
func wire(_ ...any) error { return nil }
