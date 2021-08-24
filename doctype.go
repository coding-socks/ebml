// Code generated by go run make_doctype.go. DO NOT EDIT.

package ebml

import _ "embed"

//go:embed ebml.xml
var schemaDefinition []byte

type EBML struct {
	EBMLVersion        uint
	EBMLReadVersion    uint
	EBMLMaxIDLength    uint
	EBMLMaxSizeLength  uint
	DocType            string
	DocTypeVersion     uint
	DocTypeReadVersion uint
	DocTypeExtension   []DocTypeExtension
}

type DocTypeExtension struct {
	DocTypeExtensionName    string
	DocTypeExtensionVersion uint
}
