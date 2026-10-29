//go:build tools

// Package tools pins library versions that are wired into later slices. The
// build tag keeps these imports out of the production binaries.
package tools

import (
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/riverqueue/river"
)
