// Package ebml implements a simple EBML parser.
package ebml

import (
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/vint"
	"io"
	"math/big"
	"math/bits"
)

type Element struct {
	ID       *big.Int
	DataSize *big.Int
}

// A Decoder represents an EBML parser reading a particular input stream.
type Decoder struct {
	// https://tools.ietf.org/html/rfc8794#section-11.2.4
	maxIDLength uint
	// https://tools.ietf.org/html/rfc8794#section-11.2.5
	maxSizeLength uint

	// TODO: consider using `io.ByteReader`.
	r io.ReadSeeker

	elCache *Element
}

// NewDecoder creates a new EBML parser reading from r.
func NewDecoder(r io.ReadSeeker) *Decoder {
	d := &Decoder{
		maxIDLength:   4,
		maxSizeLength: 8,
	}
	d.switchToReader(r)
	return d
}

func (d *Decoder) switchToReader(r io.ReadSeeker) {
	d.r = r
}

func (d *Decoder) skip(el *Element) error {
	_, err := d.r.Seek(el.DataSize.Int64(), io.SeekCurrent)
	return err
}

// element returns the next EBML Element in the input stream.
// At the end of the input stream, element returns nil, io.EOF.
//
// Element implements EBML specification as described by
// https://matroska-org.github.io/libebml/specs.html.
func (d *Decoder) element() (el Element, err error) {
	if d.elCache != nil {
		el = *d.elCache
		d.elCache = nil
		return el, nil
	}
	el.ID, err = d.elementID()
	if err != nil {
		return Element{}, err
	}
	el.DataSize, err = d.elementDataSize()
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

var errInvalidId = fmt.Errorf("ebml: invalid length descriptor")

// The octet length of an Element ID determines its EBML Class.
func (d *Decoder) elementID() (*big.Int, error) {
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
	return id.Val(), nil
}

func (d *Decoder) elementDataSize() (*big.Int, error) {
	b := make([]byte, d.maxSizeLength)
	// TODO: EBMLMaxSizeLength can be greater than 8
	//   https://tools.ietf.org/html/rfc8794#section-11.2.5
	if _, err := d.r.Read(b[:1]); err != nil {
		return nil, err
	}
	w := bits.LeadingZeros8(b[0]) + 1
	if _, err := d.r.Read(b[1:w]); err != nil {
		return nil, err
	}
	ds := vint.NewVint(b[:w])
	for _, b := range b[1:w] {
		if b == 255 {
			continue
		}
	}
	if allOneVint(ds.Data().Bytes(), ds.Width()) {
		return nil, nil
	}
	return ds.Data(), nil
}
