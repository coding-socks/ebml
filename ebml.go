//go:generate go run make_doctype.go

// Package ebml implements a simple EBML parser.
//
// The EBML specification is RFC 8794.
package ebml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/ebmltext"
	"github.com/coding-socks/ebml/schema"
	"io"
	"iter"
	"maps"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
)

var ErrInvalidVINTLength = ebmltext.ErrInvalidVINTWidth

var (
	docTypesMu sync.RWMutex
	docTypes   = make(map[string]*Def)

	headerDocType schema.Schema
	HeaderDef     *Def

	DefaultMaxIDLength   uint = 4
	DefaultMaxSizeLength uint = 8
)

func init() {
	err := xml.Unmarshal(schemaDefinition, &headerDocType)
	if err != nil {
		panic("cannot parse header definition")
	}
	HeaderDef, _ = NewDef(headerDocType)
}

type Def struct {
	m      map[schema.ElementID]schema.Element
	mfield map[string][]schema.Element
	Root   schema.Element
}

func NewDef(s schema.Schema) (*Def, error) {
	def := Def{
		m:      make(map[schema.ElementID]schema.Element, len(s.Elements)),
		mfield: make(map[string][]schema.Element, len(s.Elements)),
	}
	set := make(map[schema.ElementID]bool, len(s.Elements))
	var bodyRoots []schema.Element
	for _, el := range s.Elements {
		if el.Type == TypeMaster && el.Default != nil {
			return nil, fmt.Errorf("ebml: master Element %v MUST NOT declare a default value.", el.ID)
		}
		set[el.ID] = true
		def.m[el.ID] = el

		if el.Type != TypeMaster {
			i := strings.LastIndex(el.Path, "\\")
			parent := el.Path[:i]
			def.mfield[parent] = append(def.mfield[parent], el)
		}

		if strings.Count(el.Path, "\\") == 1 && el.ID != IDVoid {
			bodyRoots = append(bodyRoots, el)
		}
	}
	if len(bodyRoots) != 1 {
		return nil, errors.New("ebml: an EBML schema MUST declare exactly one EBML element at root level")
	}
	def.Root = bodyRoots[0]
	for _, el := range headerDocType.Elements {
		if set[el.ID] {
			continue
		}
		def.m[el.ID] = el
	}
	return &def, nil
}

func (d *Def) Get(id schema.ElementID) (schema.Element, bool) {
	el, ok := d.m[id]
	if !ok {
		el = UnknownSchema
	}
	return el, ok
}

func (d *Def) Fields(path string) iter.Seq[schema.Element] {
	return slices.Values(d.mfield[path])
}

func (d *Def) All() iter.Seq[schema.Element] {
	return maps.Values(d.m)
}

// Register makes a schema.Schema available by the provided doc type.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(docType string, s schema.Schema) {
	docTypesMu.Lock()
	defer docTypesMu.Unlock()
	// TODO: Validate schema
	if _, dup := docTypes[docType]; dup {
		panic("ebml: register called twice for docType " + docType)
	}
	def, err := NewDef(s)
	if err != nil {
		panic(err)
	}
	docTypes[docType] = def
}

// DocTypes returns a sorted list of the names of the registered document types.
func DocTypes() []string {
	docTypesMu.RLock()
	defer docTypesMu.RUnlock()
	list := make([]string, 0, len(docTypes))
	for name := range docTypes {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

type UnknownDocTypeError struct {
	DocType string
}

func (e UnknownDocTypeError) Error() string {
	return fmt.Sprintf("ebml: unknown DocType %q (forgotten import?)", e.DocType)
}

func Definition(docType string) (*Def, error) {
	docTypesMu.RLock()
	dt, ok := docTypes[docType]
	docTypesMu.RUnlock()
	if !ok {
		return nil, UnknownDocTypeError{DocType: docType}
	}
	return dt, nil
}

var UnknownSchema = schema.Element{
	Name:          "Unknown element",
	Documentation: []schema.Documentation{{Content: "The purpose of this object is to signal an error."}},
}

type Element struct {
	ID schema.ElementID

	// DataSize expresses the length of Element Data. Unknown data length is
	// represented with `-1`.
	//
	// With 8 octets it can have 2^56-2 possible values. That fits into int64.
	DataSize int64

	Schema schema.Element
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	r   *ebmltext.Decoder
	def *Def

	el *Element
	n  int
	// skippedErrs signals to return errors at the end of Decode.
	skippedErrs error

	window    []byte
	typeInfos map[reflect.Type]*typeInfo

	visitor Visitor
}

// NewDecoder reads and parses an EBML Document from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:   ebmltext.NewDecoder(r),
		def: HeaderDef,

		typeInfos: make(map[reflect.Type]*typeInfo),
	}
}

func (d *Decoder) SetVisitor(v Visitor) {
	d.visitor = v
}

// next reads the following element id and data size.
//
// When next encounters an ErrInvalidVINTLength or the element has UnknownSchema,
// it could be caused by damaged data or garbage in the stream. It is up
// to the caller to decide if they want to skip to the next element or
// move the reader forward by seeking one byte using io.SeekCurrent whence.
func (d *Decoder) next() (el Element, n int, err error) {
	el.ID, err = d.r.ReadElementID()
	if err != nil {
		return Element{}, n, err
	}
	n += d.r.Release()
	el.DataSize, err = d.r.ReadElementDataSize()
	if err != nil {
		return Element{}, n, err
	}
	n += d.r.Release()
	d.n = n
	sch, ok := d.def.Get(el.ID)
	if !ok {
		el.Schema = UnknownSchema
	} else {
		el.Schema = sch
	}
	if sch.Type == TypeMaster {
		d.callVisitor(el, d.r.InputOffset(), n, nil)
	}
	return el, n, err
}

var RootEl = Element{DataSize: -1, Schema: schema.Element{Name: "root", Type: TypeMaster}}

// NextOf reads the following element id and data size
// related to the given parent Element.
//
// When NextOf encounters io.EOF or end-of-element condition, it
// returns io.EOF.
//
// When NextOf encounters ErrElementOverflow fo known data size,
// you can skip the parent object, or you can read until the parent ends.
//
// See Next about ErrInvalidVINTLength.
func (d *Decoder) NextOf(parent Element, offset int64) (el Element, n int, err error) {
	if end := d.EndOfKnownDataSize(parent, offset); end {
		return Element{}, 0, io.EOF
	}
	if d.el != nil {
		el = *d.el
		d.el = nil
	} else {
		el, n, err = d.next()
		if err != nil {
			return Element{}, n, err
		}
	}
	if parent.DataSize != -1 && offset+el.DataSize > parent.DataSize {
		err = ErrElementOverflow
	}
	if end := d.EndOfUnknownDataSize(parent, el); end {
		tmp := el // This is unexpected. I cannot use the pointer to the return parameter variable.
		d.el = &tmp
		return Element{}, 0, io.EOF
	}
	return el, n, err
}

func (d *Decoder) AsSeeker() (io.Seeker, bool) {
	s, ok := d.r.AsSeeker()
	if !ok {
		return nil, false
	}
	return DecodeSeeker{d: d, ss: s}, ok
}

type DecodeSeeker struct {
	d  *Decoder
	ss io.Seeker
}

func (s DecodeSeeker) Seek(offset int64, whence int) (ret int64, err error) {
	d, ss := s.d, s.ss
	if offset != 0 && whence != io.SeekCurrent {
		d.el = nil
	}
	return ss.Seek(offset, whence)
}

// EndOfKnownDataSize tries to guess the end of an element which has a know data size.
//
// A parent with unknown data size won't raise an error but not handled as the end of the parent.
func (d *Decoder) EndOfKnownDataSize(parent Element, offset int64) bool {
	if parent.DataSize == -1 {
		return false
	}
	return offset >= parent.DataSize
}

// EndOfUnknownDataSize tries to guess the end of an element which has an unknown data size.
//
// A parent with known data size won't raise an error but not handled as the end of the parent.
func (d *Decoder) EndOfUnknownDataSize(parent Element, el Element) bool {
	if parent.DataSize != -1 {
		return false
	}
	if el.ID == IDCRC32 || el.ID == IDVoid { // global elements are child of anything
		return false
	}
	parentSch := parent.Schema
	elSch := el.Schema
	return !strings.HasPrefix(elSch.Path, parentSch.Path) || len(elSch.Path) == len(parentSch.Path)
}

type Visitor interface {
	Visit(el Element, offset int64, headerSize int, val any) (w Visitor)
}
