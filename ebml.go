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
	"github.com/coding-socks/ebml/internal/schema"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

var (
	docTypesMu sync.RWMutex
	docTypes   = make(map[string]schema.Schema)

	ebmlId  = "0x1A45DFA3"
	crc32Id = "0xBF"
	voidId  = "0xEC"

	HeaderDocType schema.Schema
)

func init() {
	err := xml.Unmarshal(schemaDefinition, &HeaderDocType)
	if err != nil {
		panic("cannot parse header definition")
	}
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
	set := make(map[string]bool, len(s.Elements))
	for i := range s.Elements {
		set[s.Elements[i].ID] = true
	}
	for i := range HeaderDocType.Elements {
		if !set[HeaderDocType.Elements[i].ID] {
			s.Elements = append(s.Elements, HeaderDocType.Elements[i])
		}
	}
	docTypes[docType] = s
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

func definition(docType string) (schema.Schema, error) {
	docTypesMu.RLock()
	dt, ok := docTypes[docType]
	docTypesMu.RUnlock()
	if !ok {
		return schema.Schema{}, fmt.Errorf("ebml: unknown docType %q (forgotten import?)", docType)
	}
	return dt, nil
}

const (
	knownDS dsMode = iota
	unknownDS
)

type dsMode int

type dataSize struct {
	m dsMode
	s int64
}

func (ds *dataSize) Known() bool {
	return ds.m == knownDS
}

func (ds *dataSize) Size() int64 {
	return ds.s
}

type Element struct {
	def schema.Element

	ID       string
	DataSize dataSize
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	Header EBML

	// https://tools.ietf.org/html/rfc8794#section-11.2.4
	maxIDLength uint
	// https://tools.ietf.org/html/rfc8794#section-11.2.5
	maxSizeLength uint

	r *Reader

	elCache *Element
}

// ReadDocument reads and parses an EBML Document from r.
func ReadDocument(r io.Reader) (*Decoder, error) {
	d := &Decoder{
		maxIDLength:   4,
		maxSizeLength: 8,
	}
	d.switchToReader(r)
	var err error
	d.Header, err = d.decodeHeader()
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Reader struct {
	pos int64
	r   io.Reader
}

func (r *Reader) Position() int64 {
	return r.pos
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.pos += int64(n)
	return n, err
}

func (d *Decoder) switchToReader(r io.Reader) {
	d.r = &Reader{r: r}
}

func (d *Decoder) skip(el Element) error {
	_, err := io.CopyN(ioutil.Discard, d.r, el.DataSize.Size())
	return err
}

type UnknownElementError struct {
	el Element
}

func (e UnknownElementError) Error() string {
	return fmt.Sprintf("ebml: unknown element: %s", e.el.ID)
}

// element returns the next EBML Element in the input stream.
// At the end of the input stream, element returns nil, io.EOF.
//
// Element implements EBML specification as described by
// https://matroska-org.github.io/libebml/specs.html.
func (d *Decoder) element(elements []schema.Element) (el Element, err error) {
	if d.elCache != nil {
		el = *d.elCache
		d.elCache = nil
	} else {
		el.ID, err = d.elementID()
		if err != nil {
			return Element{}, err
		}
		el.DataSize, err = d.elementDataSize()
		if err != nil {
			return Element{}, err
		}
	}
	var (
		found bool
		eldef schema.Element
	)
	for i := range elements {
		if elements[i].ID == el.ID {
			found = true
			eldef = elements[i]
			break
		}
	}
	if !found {
		return Element{}, &UnknownElementError{el: el}
	}
	el.def = eldef
	return el, nil
}

// https://tools.ietf.org/html/rfc8794#section-5
func validateIDData(data []byte, w int) error {
	if vintDataAllOne(data, w) {
		return errors.New("VINT_DATA MUST NOT be set to all 1")
	}
	return nil
}

var errInvalidId = fmt.Errorf("ebml: invalid length descriptor")

// The octet length of an Element ID determines its EBML Class.
func (d *Decoder) elementID() (string, error) {
	b := make([]byte, d.maxIDLength)
	// TODO: EBMLMaxIDLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.4
	if _, err := d.r.Read(b[:1]); err != nil {
		return "", err
	}
	w := vintOctetLength(b)
	if w > len(b) {
		return "", errInvalidId
	}
	if _, err := d.r.Read(b[1:w]); err != nil {
		return "", err
	}
	data := vintData(b, w)
	if err := validateIDData(data, w); err != nil {
		return "", err
	}
	return "0x" + strings.ToUpper(hex.EncodeToString(b[:w])), nil
}

func dataPad(b []byte) []byte {
	db := make([]byte, 8)
	copy(db[8-len(b):], b)
	return db
}

func (d *Decoder) elementDataSize() (dataSize, error) {
	b := make([]byte, d.maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	if _, err := d.r.Read(b[:1]); err != nil {
		return dataSize{}, err
	}
	w := vintOctetLength(b)
	if _, err := d.r.Read(b[1:w]); err != nil {
		return dataSize{}, err
	}
	ds := vintData(b, w)
	if vintDataAllOne(ds, w) {
		return dataSize{m: unknownDS}, nil
	}
	i := binary.BigEndian.Uint64(dataPad(ds))
	return dataSize{s: int64(i)}, nil
}

func DecodeDataSize(r io.Reader) (dataSize, error) {
	b := make([]byte, 8)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	if _, err := r.Read(b[:1]); err != nil {
		return dataSize{}, err
	}
	w := vintOctetLength(b)
	if _, err := r.Read(b[1:w]); err != nil {
		return dataSize{}, err
	}
	ds := vintData(b, w)
	if vintDataAllOne(ds, w) {
		return dataSize{m: unknownDS}, nil
	}
	i := binary.BigEndian.Uint64(dataPad(ds))
	return dataSize{s: int64(i)}, nil
}
