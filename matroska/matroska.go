//go:generate go run make_doctype.go

// Package matroska contains types and structures for parsing
// matroska (.mkv, .mk3d, .mka, .mks) files.
package matroska

import "github.com/coding-socks/ebml"

func init() {
	ebml.Register("matroska", &DocType)
}
