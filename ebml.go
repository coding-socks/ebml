//go:generate go run make_doctype.go

// Package ebml implements a simple EBML parser.
//
// The EBML specification is RFC 8794.
package ebml

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

var (
	docTypesMu sync.RWMutex
	docTypes   = make(map[string]Definition)

	CRC32 = NewDefinition("BF", TypeBinary, "CRC-32", nil, nil)
	Void  = NewDefinition("EC", TypeBinary, "CRC-32", nil, nil)
)

// Register makes a DocType available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, docType *Definition) {
	docTypesMu.Lock()
	defer docTypesMu.Unlock()
	if docType == nil {
		panic("ebml: Register docType is nil")
	}
	// TODO: Validate schema
	if _, dup := docTypes[name]; dup {
		panic("ebml: Register called twice for docType " + name)
	}
	docTypes[name] = *docType
}

// Drivers returns a sorted list of the names of the registered drivers.
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

func getDefinition(docType string) (Definition, error) {
	docTypesMu.RLock()
	dt, ok := docTypes[docType]
	docTypesMu.RUnlock()
	if !ok {
		return Definition{}, fmt.Errorf("ebml: unknown docType %q (forgotten import?)", docType)
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
	ID       string
	DataSize dataSize

	Definition Definition
}

// A Decoder represents an EBML parser reading a particular input stream.
// Decoding one document consists of calling DecodeHeader and then DecodeBody
// in that order.
type Decoder struct {
	// https://tools.ietf.org/html/rfc8794#section-11.2.4
	maxIDLength uint
	// https://tools.ietf.org/html/rfc8794#section-11.2.5
	maxSizeLength uint

	r *Reader

	elCache *Element

	headerDefinition Definition
	bodyDefinition   Definition
}

// NewDecoder creates a new EBML parser reading from r.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		maxIDLength:   4,
		maxSizeLength: 8,
	}
	d.switchToReader(r)
	return d
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

func (d *Decoder) skip(el *Element) error {
	_, err := io.CopyN(ioutil.Discard, d.r, el.DataSize.Size())
	return err
}

type UnknownElementError struct {
	el Element
}

func (e UnknownElementError) Error() string {
	return fmt.Sprintf("ebml: unknown element: 0x%s", e.el.ID)
}

// element returns the next EBML Element in the input stream.
// At the end of the input stream, element returns nil, io.EOF.
//
// Element implements EBML specification as described by
// https://matroska-org.github.io/libebml/specs.html.
func (d *Decoder) element(defs []Definition) (el Element, err error) {
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
		eldef Definition
	)
	for i := range defs {
		def := defs[i]
		if def.ID == el.ID {
			found = true
			eldef = def
			break
		}
	}
	if !found {
		switch {
		default:
			return Element{}, &UnknownElementError{el: el}
		case CRC32.ID == el.ID:
			eldef = CRC32
		case Void.ID == el.ID:
			eldef = Void
		}
	}
	el.Definition = eldef
	return el, nil
}

// https://tools.ietf.org/html/rfc8794#section-5
func validateIDData(data []byte, w int) error {
	if vintDataAllZero(data, w) {
		return errors.New("VINT_DATA MUST NOT be set to all 0")
	}
	if vintDataAllOne(data, w) {
		return errors.New("VINT_DATA MUST NOT be set to all 1")
	}
	if vintDataAvailableShorter(data, w) {
		return errors.New("a shorter VINT_DATA encoding is available")
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
	return hex.EncodeToString(b[:w]), nil
}

func dataPad(b []byte) []byte {
	db := make([]byte, 8)
	for i, j := 8-len(b), 0; i < 8; i, j = i+1, j+1 {
		db[i] = b[j]
	}
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

type Definition struct {
	ID       string
	Type     string
	Name     string
	Default  interface{}
	Children []Definition
}

func NewDefinition(id string, t, name string, def interface{}, children []Definition) Definition {
	id = strings.TrimPrefix(id, "0x")
	_, err := hex.DecodeString(id)
	if err != nil {
		panic(err)
	}
	return Definition{ID: strings.ToLower(id), Type: t, Name: name, Default: def, Children: children}
}
