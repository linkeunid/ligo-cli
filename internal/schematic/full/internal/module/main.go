package module

import "github.com/linkeunid/ligo"

// Main is the root module that composes all application modules.
func MainModule() ligo.Module {
	return ligo.NewModule("main",
		ligo.Imports(
			AuthModule(),
			UserModule(),
			FileModule(),
			HealthModule(),
			RootModule(),
		),
	)
}
