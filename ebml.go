// Package ebml implements a simple EBML parser.
package ebml

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/vint"
	"io"
	"math/big"
	"math/bits"
)

type Element struct {
	ID       *vint.Vint
	DataSize *vint.Vint
	Data     io.Reader
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	// https://tools.ietf.org/html/rfc8794#section-11.2.4
	maxIDLength uint
	// https://tools.ietf.org/html/rfc8794#section-11.2.5
	maxSizeLength uint

	// TODO: consider using `io.ByteReader`.
	r io.ByteReader
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

func (d *Decoder) switchToReader(r io.Reader) {
	// Get efficient byte at a time reader.
	// Assume that if reader has its own
	// ReadByte, it's efficient enough.
	// Otherwise, use bufio.
	if rb, ok := r.(io.ByteReader); ok {
		d.r = rb
	} else {
		d.r = bufio.NewReader(r)
	}
}

// Element returns the next EBML element in the input stream.
// At the end of the input stream, Element returns nil, io.EOF.
//
// Element implements EBML specification as described by
// https://matroska-org.github.io/libebml/specs.html.
func (d *Decoder) Element() (el Element, err error) {
	el.ID, err = d.elementID()
	if err != nil {
		return Element{}, err
	}
	el.DataSize, err = d.elementDataSize()
	if err != nil {
		return Element{}, err
	}
	el.Data, err = d.elementData(el.DataSize)
	if err != nil {
		return Element{}, err
	}
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

// The octet length of an Element ID determines its EBML Class.
func (d *Decoder) elementID() (*vint.Vint, error) {
	var err error
	b := make([]byte, d.maxIDLength)
	// TODO: EBMLMaxIDLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.4
	if b[0], err = d.r.ReadByte(); err != nil {
		return nil, err
	}
	w := bits.LeadingZeros8(b[0]) + 1
	if w > len(b) {
		return nil, fmt.Errorf("ebml: invalid length descriptor: %08b", w)
	}
	for i := 1; i < w; i++ {
		if b[i], err = d.r.ReadByte(); err != nil {
			return nil, err
		}
	}
	id := vint.NewVint(b[:w])
	if err := validateID(id); err != nil {
		return nil, err
	}
	return id, nil
}

func (d *Decoder) elementDataSize() (*vint.Vint, error) {
	var err error
	b := make([]byte, d.maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	if b[0], err = d.r.ReadByte(); err != nil {
		return nil, err
	}
	w := bits.LeadingZeros8(b[0]) + 1
	for i := 1; i < w; i++ {
		if b[i], err = d.r.ReadByte(); err != nil {
			return nil, err
		}
	}
	ds := vint.NewVint(b[:w])
	for _, b := range b[1:w] {
		if b == 255 {
			continue
		}
	}
	if allOneVint(ds.Data().Bytes(), ds.Width()) {
		return vint.Inf, nil
	}
	return ds, nil
}

func (d *Decoder) elementData(ds *vint.Vint) (io.Reader, error) {
	if ds.IsInf() {
		// TODO: Handle unknown data size
		//  https://tools.ietf.org/html/rfc8794#section-6.2
		panic("ebml: Unknown data size is not implemented")
	}
	var buf bytes.Buffer
	dataSize := new(big.Int).Set(ds.Data())
	for dataSize.Cmp(big.NewInt(0)) > 0 {
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		buf.WriteByte(b)
		dataSize = dataSize.Sub(dataSize, big.NewInt(1))
	}
	return &buf, nil
}
