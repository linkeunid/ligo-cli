// Package ext provides templates for generating Ligo extension modules.
package ext

import "embed"

//go:embed *.tmpl
var FS embed.FS
