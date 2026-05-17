// Package ext provides templates for generating Ligo extension modules.
package ext

import "embed"

// Embed includes both top-level .tmpl files and the dotted CI / lint
// config templates (nested under .github/). The walker strips the
// .tmpl suffix before writing each file to disk.
//
//go:embed *.tmpl
//go:embed .golangci.yml.tmpl
//go:embed .github
var FS embed.FS
