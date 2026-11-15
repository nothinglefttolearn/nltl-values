// Package nltlvalues embeds the raw value files (tokens/ and shapes/) so Go
// consumers get them at the pinned module version.
//
// Most code should use the typed accessors in github.com/nothinglefttolearn/nltl-values/go
// instead. This file lives at the repo root only because go:embed cannot
// reach parent directories.
package nltlvalues

import "embed"

// FS holds tokens/*.json and shapes/*.json, at the same paths as in the repo.
//
//go:embed tokens/*.json shapes/*.json
var FS embed.FS
