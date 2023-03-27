//go:generate go run make_doctype.go

// Package ebml implements a simple EBML parser.
//
// The EBML specification is RFC 8794.
package ebml

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/schema"
	"io"
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
	m    map[string]schema.Element
	Root schema.Element
}

func NewDef(s schema.Schema) (*Def, error) {
	def := Def{
		m: make(map[string]schema.Element, len(s.Elements)),
	}
	set := make(map[string]bool, len(s.Elements))
	for _, el := range s.Elements {
		if el.Type == TypeMaster && el.Default != nil {
			return nil, fmt.Errorf("ebml: master Element %s MUST NOT declare a default value.", el.ID)
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

func (d *Def) Get(id string) (schema.Element, bool) {
	el, ok := d.m[id]
	return el, ok
}

func (d *Def) Values() []schema.Element {
	els := make([]schema.Element, len(d.m))
	var i int
	for s := range d.m {
		els[i] = d.m[s]
		i++
	}
	return els
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

const (
	unknownDS dsMode = iota
	knownDS
)

type dsMode int

type DataSize struct {
	m dsMode
	s int64
}

func (ds *DataSize) Known() bool {
	return ds.m == knownDS
}

func (ds *DataSize) Size() int64 {
	return ds.s
}

type Element struct {
	ID       string
	DataSize DataSize
}

// Reader provides a low level API to interacts with EBML documents.
// Use directly with caution.
type Reader struct {
	r    io.ReaderAt
	base int64
	off  int64

	// https://tools.ietf.org/html/rfc8794#section-11.2.4
	MaxIDLength uint
	// https://tools.ietf.org/html/rfc8794#section-11.2.5
	MaxSizeLength uint
}

func NewReader(r io.ReaderAt) *Reader {
	return &Reader{
		r: r,

		MaxIDLength:   DefaultMaxIDLength,
		MaxSizeLength: DefaultMaxSizeLength,
	}
}

func (r *Reader) Next() (el Element, n int, err error) {
	el.ID, n, err = ReadElementID(r, r.MaxIDLength)
	if err != nil {
		return Element{}, n, err
	}
	var m int
	el.DataSize, m, err = ReadElementDataSize(r, r.MaxSizeLength)
	n += m
	if err != nil {
		return Element{}, n, err
	}
	return el, n, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.ReadAt(p, r.off)
	r.off += int64(n)
	return
}

var errWhence = errors.New("Seek: invalid whence")
var errOffset = errors.New("Seek: invalid offset")

func (r *Reader) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	default:
		return 0, errWhence
	case io.SeekStart:
		offset += r.base
	case io.SeekCurrent:
		offset += r.off
	case io.SeekEnd:
		// TODO: not sure how to handle this
		panic("ebml: not able to seek relative to the end")
	}
	if offset < r.base {
		return 0, errOffset
	}
	r.off = offset
	return offset - r.base, nil
}

func (r *Reader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, io.EOF
	}
	off += r.base
	return r.r.ReadAt(p, off)
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	r   *Reader
	def *Def

	el *Element
}

// NewDecoder reads and parses an EBML Document from r.
func NewDecoder(r io.ReaderAt) *Decoder {
	return &Decoder{
		r:   NewReader(r),
		def: HeaderDef,
	}
}

func (d *Decoder) Next() (el Element, n int, err error) {
	el, i, err := d.r.Next()
	d.el = &el
	return el, i, err
}

func (d *Decoder) Seek(offset int64, whence int) (ret int64, err error) {
	d.el = nil
	return d.r.Seek(offset, whence)
}

type UnknownDefinitionError struct {
	id string
}

func (u UnknownDefinitionError) ID() string {
	return u.id
}

func (u UnknownDefinitionError) Error() string {
	return fmt.Sprintf("ebml: element definition not found for %s", u.id)
}

// EndOfElement tries to guess the end of an element.
//
// Offset is ignored when element has unknown size.
func (d *Decoder) EndOfElement(parent Element, el Element, offset int64) (bool, error) {
	if parent.DataSize.Known() {
		if offset > parent.DataSize.Size() {
			return true, ErrElementOverflow
		}
		return offset == parent.DataSize.Size(), nil
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
	return fmt.Sprintf("ebml: unknown element: %s", e.el.ID)
}

var ErrInvalidVINTLength = fmt.Errorf("ebml: invalid length descriptor")

// ReadElementID reads an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-5
func ReadElementID(r io.Reader, maxIDLength uint) (id string, n int, err error) {
	b := make([]byte, maxIDLength)
	// TODO: EBMLMaxIDLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.4
	n, err = r.Read(b[:1])
	if err != nil {
		return "", n, err
	}
	w := vintOctetLength(b)
	if w > len(b) {
		return "", 1, ErrInvalidVINTLength
	}
	if w > 1 {
		m, err := r.Read(b[1:w])
		n += m
		if err != nil {
			return "", n, err
		}
	}
	data := vintData(b, w)
	if vintDataAllOne(data, w) {
		return "", n, errors.New("VINT_DATA MUST NOT be set to all 1")
	}
	return "0x" + strings.ToUpper(hex.EncodeToString(b[:w])), n, nil
}

func dataPad(b []byte) []byte {
	db := make([]byte, 8)
	copy(db[8-len(b):], b)
	return db
}

// ReadElementDataSize reads an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-6
func ReadElementDataSize(r io.Reader, maxSizeLength uint) (ds DataSize, n int, err error) {
	b := make([]byte, maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	n, err = r.Read(b[:1])
	if err != nil {
		return DataSize{}, n, err
	}
	w := vintOctetLength(b)
	m, err := r.Read(b[1:w])
	n += m
	if err != nil {
		return DataSize{}, n, err
	}
	d := vintData(b, w)
	if vintDataAllOne(d, w) {
		return DataSize{}, n, nil
	}
	i := binary.BigEndian.Uint64(dataPad(d))
	return DataSize{m: knownDS, s: int64(i)}, n, nil
}
