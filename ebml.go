//go:generate go run make_doctype.go

// Package ebml implements a simple EBML parser.
//
// The EBML specification is RFC 8794.
package ebml

import (
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/ebmltext"
	"github.com/coding-socks/ebml/schema"
	"io"
	"iter"
	"maps"
	"reflect"
	"sort"
	"strings"
	"sync"
)

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
	m    map[schema.ElementID]schema.Element
	Root schema.Element
}

func NewDef(s schema.Schema) (*Def, error) {
	def := Def{
		m: make(map[schema.ElementID]schema.Element, len(s.Elements)),
	}
	set := make(map[schema.ElementID]bool, len(s.Elements))
	for _, el := range s.Elements {
		if el.Type == TypeMaster && el.Default != nil {
			return nil, fmt.Errorf("ebml: master Element %v MUST NOT declare a default value.", el.ID)
		}
		set[el.ID] = true
		def.m[el.ID] = el
	}
	var bodyRoots []schema.Element
	for _, el := range def.m {
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
	return el, ok
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

type Element struct {
	ID schema.ElementID

	// DataSize expresses the length of Element Data. Unknown data length is
	// represented with `-1`.
	//
	// With 8 octets it can have 2^56-2 possible values. That fits into int64.
	DataSize int64
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	r   *ebmltext.Decoder
	def *Def

	el *Element
	// elOverflow signals to return ErrElementOverflow at the end of decode.
	elOverflow bool

	window    []byte
	typeInfos map[reflect.Type]*typeInfo
}

// NewDecoder reads and parses an EBML Document from r.
func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{
		r:   ebmltext.NewDecoder(r),
		def: HeaderDef,

		typeInfos: make(map[reflect.Type]*typeInfo),
	}
}

// Next reads the following element id and data size.
// It must be called before Decode.
func (d *Decoder) Next() (el Element, n int, err error) {
	start := d.r.InputOffset()
	el.ID, err = d.r.ReadElementID()
	if err != nil {
		return Element{}, int(d.r.InputOffset() - start), err
	}
	d.r.Release()
	el.DataSize, err = d.r.ReadElementDataSize()
	if err != nil {
		return Element{}, int(d.r.InputOffset() - start), err
	}
	d.r.Release()
	d.el = &el
	return el, int(d.r.InputOffset() - start), err
}

// NextOf reads the following element id and data size
// related to the given parent Element.
//
// When NextOf encounters an error or end-of-element condition it
// return EOE error.
func (d *Decoder) NextOf(parent Element, offset int64) (el Element, n int, err error) {
	if end, err := d.EndOfKnownDataSize(parent, offset); err != nil {
		return Element{}, 0, err
	} else if end {
		return Element{}, 0, io.EOF
	}
	el, n, err = d.Next()
	if err != nil {
		return Element{}, n, err
	}
	if end, err := d.EndOfUnknownDataSize(parent, el); err != nil {
		return Element{}, n, err
	} else if end {
		d.r.Seek(int64(-n), io.SeekCurrent)
		return Element{}, 0, io.EOF
	}
	return el, n, nil
}

func (d *Decoder) Seek(offset int64, whence int) (ret int64, err error) {
	if offset != 0 && whence != io.SeekCurrent {
		d.el = nil
	}
	return d.r.Seek(offset, whence)
}

type UnknownDefinitionError struct {
	id schema.ElementID
}

func (u UnknownDefinitionError) ID() schema.ElementID {
	return u.id
}

func (u UnknownDefinitionError) Error() string {
	return fmt.Sprintf("ebml: element definition not found for %v", u.id)
}

// EndOfKnownDataSize tries to guess the end of an element which has a know data size.
//
// A parent with unknown data size won't raise an error but not handled as the end of the parent.
func (d *Decoder) EndOfKnownDataSize(parent Element, offset int64) (bool, error) {
	if parent.DataSize == -1 {
		return false, nil
	}
	if offset > parent.DataSize {
		return true, ErrElementOverflow
	}
	return offset == parent.DataSize, nil
}

// EndOfUnknownDataSize tries to guess the end of an element which has an unknown data size.
//
// A parent with known data size won't raise an error but not handled as the end of the parent.
func (d *Decoder) EndOfUnknownDataSize(parent Element, el Element) (bool, error) {
	if parent.DataSize != -1 {
		return false, nil
	}
	if el.ID == IDCRC32 || el.ID == IDVoid { // global elements are child of anything
		return false, nil
	}
	def, ok := d.def.Get(parent.ID)
	if !ok {
		return false, &UnknownDefinitionError{parent.ID}
	}
	nextDef, ok := d.def.Get(el.ID)
	if !ok {
		return false, &UnknownDefinitionError{el.ID}
	}
	return !strings.HasPrefix(nextDef.Path, def.Path) || len(nextDef.Path) == len(def.Path), nil
}

type UnknownElementError struct {
	el Element
}

func (e UnknownElementError) Error() string {
	return fmt.Sprintf("ebml: unknown element: %v", e.el.ID)
}

var ErrInvalidVINTLength = fmt.Errorf("ebml: invalid length descriptor")

// ReadElementID reads an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-5
func ReadElementID(r io.Reader, maxIDLength uint) (id schema.ElementID, n int, err error) {
	b := make([]byte, maxIDLength)
	// TODO: EBMLMaxIDLength can be greater than 8
	//   https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.4
	n, err = r.Read(b[:1])
	if err != nil {
		return 0, n, err
	}
	w := vintOctetLength(b)
	if w > len(b) {
		return 0, 1, ErrInvalidVINTLength
	}
	if w > 1 {
		m, err := r.Read(b[1:w])
		n += m
		if err != nil {
			return 0, n, err
		}
	}
	data := vintData(b, w)
	if vintDataAllOne(data, w) {
		return 0, n, errors.New("VINT_DATA MUST NOT be set to all 1")
	}
	i := binary.BigEndian.Uint64(dataPad(b[:w]))
	return schema.ElementID(i), n, nil
}

func dataPad(b []byte) []byte {
	db := make([]byte, 8)
	copy(db[8-len(b):], b)
	return db
}

// ReadElementDataSize reads an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-6
func ReadElementDataSize(r io.Reader, maxSizeLength uint) (ds int64, n int, err error) {
	b := make([]byte, maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.5
	n, err = r.Read(b[:1])
	if err != nil {
		return 0, n, err
	}
	w := vintOctetLength(b)
	m, err := r.Read(b[1:w])
	n += m
	if err != nil {
		return 0, n, err
	}
	d := vintData(b, w)
	if vintDataAllOne(d, w) {
		return -1, n, nil
	}
	i := binary.BigEndian.Uint64(dataPad(d))
	return int64(i), n, nil
}
