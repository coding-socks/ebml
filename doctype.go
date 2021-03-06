// Code generated by go run make_doctype.go. DO NOT EDIT.

package ebml

var HeaderDocType = NewDefinition("0x1A45DFA3", TypeMaster, "EBML", nil, []Definition{
	NewDefinition("0x4286", TypeUinteger, "EBMLVersion", uint(1), nil),
	NewDefinition("0x42F7", TypeUinteger, "EBMLReadVersion", uint(1), nil),
	NewDefinition("0x42F2", TypeUinteger, "EBMLMaxIDLength", uint(4), nil),
	NewDefinition("0x42F3", TypeUinteger, "EBMLMaxSizeLength", uint(8), nil),
	NewDefinition("0x4282", TypeString, "DocType", "ebml", nil),
	NewDefinition("0x4287", TypeUinteger, "DocTypeVersion", uint(1), nil),
	NewDefinition("0x4285", TypeUinteger, "DocTypeReadVersion", uint(1), nil),
	NewDefinition("0x4281", TypeMaster, "DocTypeExtension", nil, []Definition{
		NewDefinition("0x4283", TypeString, "DocTypeExtensionName", nil, nil),
		NewDefinition("0x4284", TypeUinteger, "DocTypeExtensionVersion", nil, nil),
	}),
})

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
