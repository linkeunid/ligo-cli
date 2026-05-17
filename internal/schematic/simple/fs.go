package simple

import "embed"

// Embed includes regular files and the dotted CI / lint config files
// that Go's default //go:embed glob would otherwise skip.
//
//go:embed *
//go:embed .golangci.yml.tmpl
//go:embed .github
var FS embed.FS
