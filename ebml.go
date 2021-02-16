//go:generate go run make_definition.go

// Package ebml implements a simple EBML parser.
package ebml

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/vint"
	"io"
	"io/ioutil"
	"math/bits"
	"sort"
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
	ID       []byte
	DataSize dataSize

	Definition Definition
}

func (e Element) HexID() string {
	return hex.EncodeToString(e.ID)
}

// A Decoder represents an EBML parser reading a particular input stream.
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

// NewDecoder creates a new EBML parser reading from r.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		maxIDLength:   4,
		maxSizeLength: 8,

		headerDefinition: headerDefinition,
	}
	d.switchToReader(r)
	return d
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
	return fmt.Sprintf("ebml: unknown element: 0x%x", e.el.ID)
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
		def   Definition
	)
	for i := range defs {
		definition := defs[i]
		if bytes.Compare(definition.ID, el.ID) == 0 {
			found = true
			def = definition
			break
		}
	}
	if !found {
		switch {
		default:
			return Element{}, &UnknownElementError{el: el}
		case bytes.Compare(CRC32.ID, el.ID) == 0:
			def = CRC32
		case bytes.Compare(Void.ID, el.ID) == 0:
			def = Void
		}
	}
	el.Definition = def
	return el, nil
}

// https://tools.ietf.org/html/rfc8794#section-5
func validateID(id *vint.Vint) error {
	b := id.Data().Bytes()
	if len(b) == 0 {
		return errors.New("VINT_DATA MUST NOT be set to all 0")
	}
	if allOneVint(b, id.Width()) {
		return errors.New("VINT_DATA MUST NOT be set to all 1")
	}
	if shorterAvailableVint(b, id.Width()) {
		return errors.New("a shorter VINT_DATA encoding is available")
	}
	return nil
}

var errInvalidId = fmt.Errorf("ebml: invalid length descriptor")

// The octet length of an Element ID determines its EBML Class.
func (d *Decoder) elementID() ([]byte, error) {
	b := make([]byte, d.maxIDLength)
	// TODO: EBMLMaxIDLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.4
	if _, err := d.r.Read(b[:1]); err != nil {
		return nil, err
	}
	w := bits.LeadingZeros8(b[0]) + 1
	if w > len(b) {
		return nil, errInvalidId
	}
	if _, err := d.r.Read(b[1:w]); err != nil {
		return nil, err
	}
	id := vint.NewVint(b[:w])
	if err := validateID(id); err != nil {
		return nil, err
	}
	return id.Val().Bytes(), nil
}

func (d *Decoder) elementDataSize() (dataSize, error) {
	b := make([]byte, d.maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	if _, err := d.r.Read(b[:1]); err != nil {
		return dataSize{}, err
	}
	w := bits.LeadingZeros8(b[0]) + 1
	if _, err := d.r.Read(b[1:w]); err != nil {
		return dataSize{}, err
	}
	ds := vint.NewVint(b[:w])
	for _, b := range b[1:w] {
		if b == 255 {
			continue
		}
	}
	if allOneVint(ds.Data().Bytes(), ds.Width()) {
		return dataSize{m: unknownDS}, nil
	}
	return dataSize{s: ds.Data().Int64()}, nil
}
