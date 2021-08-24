// Package schema contains structs for reading xml definitions
// for ebml schema.
//
// TODO: Generated from EBMLSchema.xsd.
//
// Internal use only.
package schema

import (
	"encoding/xml"
	"github.com/coding-socks/ebml/internal/ebmlpath"
	"reflect"
	"strconv"
)

var (
	PurposeDefinition = "definition"
	PurposeRationale  = "rationale"
	PurposeReferences = "references"
	PurposeUsageNotes = "usage notes"
)

type Documentation struct {
	Content string `xml:",chardata"`
	Lang    string `xml:"lang,attr"`
	Purpose string `xml:"purpose"`
}

var (
	NoteAttributeMinOccurs = "minOccurs"
	NoteAttributeMaxOccurs = "maxOccurs"
	NoteAttributeRange     = "range"
	NoteAttributeLength    = "length"
	NoteAttributeDefault   = "default"
	NoteAttributeMinver    = "minver"
	NoteAttributeMaxver    = "maxver"
)

type Note struct {
	Content       string `xml:",chardata"`
	NoteAttribute string `xml:"note_attribute,attr"`
}

type Enum struct {
	Documentation []Documentation `xml:"documentation"`
	Label         string          `xml:"label,attr"`
	Value         string          `xml:"value,attr"`
}

type Restriction struct {
	Enum []Enum `xml:"enum"`
}

type Extension struct {
	Type       string     `xml:"type,attr"`
	Attributes []xml.Attr `xml:",any,attr"`
}

var (
	TypeInteger  = "integer"
	TypeUinteger = "uinteger"
	TypeFloat    = "float"
	TypeString   = "string"
	TypeDate     = "date"
	TypeUtf8     = "utf-8"
	TypeMaster   = "master"
	TypeBinary   = "binary"
)

type Element struct {
	Documentation      []Documentation `xml:"documentation"`
	ImplementationNote []Note          `xml:"implementation_note"`
	Restriction        *Restriction    `xml:"restriction"`
	Extension          []Extension     `xml:"extension"`

	Name               string       `xml:"name,attr"`
	Path               Path         `xml:"path,attr"`
	ID                 string       `xml:"id,attr"`
	MinOccurs          int          `xml:"minOccurs,attr"`
	MaxOccurs          UnboundedInt `xml:"maxOccurs,attr"`
	Range              string       `xml:"range,attr"`
	Length             string       `xml:"length,attr"`
	Default            *string      `xml:"default,attr"`
	Type               string       `xml:"type,attr"`
	UnknownSizeAllowed bool         `xml:"unknownsizeallowed,attr"`
	Recursive          bool         `xml:"recursive,attr"`
	Recurring          bool         `xml:"recurring,attr"`
	MinVer             int          `xml:"minver,attr"`
	MaxVer             int          `xml:"maxver,attr"`
}

type Path struct {
	*ebmlpath.PathExp
}

func (p *Path) UnmarshalXMLAttr(attr xml.Attr) error {
	var err error
	p.PathExp, err = ebmlpath.Compile(attr.Value)
	return err
}

type UnboundedInt struct {
	unbounded bool
	val       int
}

func (u UnboundedInt) Unbounded() bool {
	return u.unbounded
}

func (u UnboundedInt) Val() int {
	return u.val
}

func (u *UnboundedInt) UnmarshalXMLAttr(attr xml.Attr) error {
	if attr.Value == "unbounded" {
		*u = UnboundedInt{unbounded: true}
		return nil
	}
	i, err := strconv.ParseInt(attr.Value, 10, 64)
	if err != nil {
		return err
	}
	*u = UnboundedInt{val: int(i)}
	return nil
}

func (s *Element) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type fw Element // prevent recursion
	item := fw{
		MinOccurs: 0, // default="0"
		// TODO: consider using pointer to differentiate between not set and real value
		// [...] If the maxOccurs attribute is not present, then there is no
		// upper bound for the permitted number of occurrences [...]
		// https://www.rfc-editor.org/rfc/rfc8794#name-maxoccurs
		MaxOccurs:          UnboundedInt{unbounded: true}, // default="unbounded"
		UnknownSizeAllowed: false,                         // default="false"
		Recursive:          false,                         // default="false"
		Recurring:          false,                         // default="false"
		MinVer:             1,                             // default="1"
	}
	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}
	*s = (Element)(item)
	return nil
}

type Schema struct {
	Elements []Element `xml:"element"`

	DocType string `xml:"docType,attr"`
	Version int    `xml:"version,attr"`
	EBML    uint   `xml:"ebml,attr"`
}

// https://stackoverflow.com/a/26957888/2231168
func (s *Schema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type fw Schema // prevent recursion
	item := fw{
		EBML: 1, // default="1"
	}
	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}
	*s = (Schema)(item)
	return nil
}

// Query returns the first element which matches the given path.
func (s *Schema) Query(path string) (Element, bool) {
	for i := range s.Elements {
		if s.Elements[i].Path.Match(path) {
			return s.Elements[i], true
		}
	}
	return Element{}, false
}

// QueryChildren returns all elements which can be a direct children
// of the given path.
func (s *Schema) QueryChildren(path string) []Element {
	var e []Element
	for i := range s.Elements {
		el := s.Elements[i]
		p := ebmlpath.Join(path, el.Name)
		if el.Path.Match(p) {
			e = append(e, el)
		}
	}
	return e
}

type TreeNode struct {
	El       Element
	children map[string]*TreeNode
	order    []string
}

func NewTreeNode(el Element) *TreeNode {
	return &TreeNode{El: el, children: make(map[string]*TreeNode)}
}

func (n *TreeNode) Put(key string, el *TreeNode) {
	if _, ok := n.children[key]; !ok {
		n.order = append(n.order, key)
	}
	n.children[key] = el
}

func (n *TreeNode) Get(key string) *TreeNode {
	return n.children[key]
}

func (n *TreeNode) VisitAll(f func(node *TreeNode)) {
	for _, key := range n.order {
		f(n.children[key])
	}
}

func ResolveGoType(s, name string) string {
	switch s {
	case TypeInteger:
		return reflect.Int.String()
	case TypeUinteger:
		return reflect.Uint.String()
	case TypeFloat:
		return reflect.Float64.String()
	case TypeString:
		// TODO: Enforce ASCII only characters (in the range of 0x20 to 0x7E).
		//  https://www.rfc-editor.org/rfc/rfc8794#name-string-element
		return reflect.String.String()
	case TypeDate:
		return "time.Time"
	case TypeUtf8:
		return reflect.String.String()
	case TypeMaster:
		return name
	case TypeBinary:
		return "[]byte"
	}
	return s
}
