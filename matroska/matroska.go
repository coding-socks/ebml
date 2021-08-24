//go:generate go run make_doctype.go

// Package matroska contains types and structures for parsing
// matroska (.mkv, .mk3d, .mka, .mks) files.
package matroska

import (
	"encoding/xml"
	"github.com/coding-socks/ebml"
	"github.com/coding-socks/ebml/internal/schema"
)

func init() {
	var s schema.Schema
	if err := xml.Unmarshal(DocType, &s); err != nil {
		panic("not able to parse matroska schema: " + err.Error())
	}

	ebml.Register(s.DocType, s)
}
