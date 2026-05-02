package module

import "github.com/linkeunid/ligo"

// Main is the root module that composes all application modules.
func Main() ligo.Module {
	return ligo.NewModule("main",
		ligo.Imports(
			Auth(),
			User(),
			File(),
			Health(),
			Root(),
		),
	)
}
