package ebmltext

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/schema"
	"io"
	"math"
	"time"
)

var thirdMillennium = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)

type Token struct {
	id       schema.ElementID
	dataSize int64

	start int64
	end   int64
}

type Decoder struct {
	// https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.4
	MaxIDLength uint
	// https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.5
	MaxSizeLength uint
	offset        int64
	releasable    int

	r byteReader
}

func NewDecoder(r io.ReadSeeker) *Decoder {
	return &Decoder{
		MaxIDLength:   8,
		MaxSizeLength: 8,

		r: byteReader{r: r},
	}
}

func (d *Decoder) InputOffset() int64 {
	return d.offset
}

var (
	ErrInvalidVINTWidth = fmt.Errorf("ebmltext: invalid VINT_WIDTH")
	ErrAllOneVINT       = fmt.Errorf("ebmltext: VINT_DATA MUST NOT be set to all 1")
)

// ReadElementID reads an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-5
//
// TODO: EBMLMaxIDLength can be greater than 8 https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.4
func (d *Decoder) ReadElementID() (id schema.ElementID, err error) {
	d.r.extend()
	window := d.r.window()
	if len(window) == 0 {
		return 0, io.EOF
	}
	vint, w, err := ReadVint(window[:d.MaxIDLength])
	if err != nil {
		return 0, ErrInvalidVINTWidth
	}
	if vintAllOne(vint) {
		return 0, ErrAllOneVINT
	}
	d.releasable = w
	return schema.ElementID(vint), nil
}

// ReadElementDataSize reads an Element Data Size based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-6
//
// TODO: EBMLMaxSizeLength can be greater than 8 https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.5
func (d *Decoder) ReadElementDataSize() (ds int64, err error) {
	d.r.extend()
	window := d.r.window()
	if len(window) == 0 {
		return 0, io.EOF
	}
	vintD, w, err := ReadVintData(window[:d.MaxSizeLength])
	if err != nil {
		return 0, err
	}
	if w > int(d.MaxSizeLength) {
		return 0, ErrInvalidVINTWidth
	}
	d.releasable = w
	if vintDataAllOne(vintD, w) {
		return -1, nil
	}
	return int64(vintD), nil
}

// Release releases n bytes from the internal buffer after reading
// an Element ID or an Element Data Size, where n is the number
// of bytes read by these operations.
//
// This is useful if the file is damaged and have to look for
// valid Elements.
func (d *Decoder) Release() int {
	d.offset += int64(d.releasable)
	d.r.release(d.releasable)
	return d.releasable
}

func (d *Decoder) Read(b []byte) (int, error) {
	n, err := d.r.Read(b)
	d.offset += int64(n)
	return n, err
}

func (d *Decoder) Seek(offset int64, whence int) (int64, error) {
	start, _ := d.r.Seek(0, io.SeekCurrent)
	end, err := d.r.Seek(offset, whence)
	d.offset += end - start
	return end, err
}

// Int reads an int64 based on https://www.rfc-editor.org/rfc/rfc8794.html#section-7.1
func Int(b []byte) (int64, error) {
	if len(b) > 8 {
		return 0, errors.New("ebml: max length for an unsigned integer is eight octets")
	}
	i := int64(0)
	for _, bb := range b {
		i = (i << 8) | int64(bb)
	}
	return i, nil
}

// Uint reads an uint64 based on https://www.rfc-editor.org/rfc/rfc8794.html#section-7.2
func Uint(b []byte) (uint64, error) {
	if len(b) > 8 {
		return 0, errors.New("ebml: max length for an unsigned integer is eight octets")
	}
	i := uint64(0)
	for _, bb := range b {
		i = (i << 8) | uint64(bb)
	}
	return i, nil
}

// Float reads a float64 based on https://www.rfc-editor.org/rfc/rfc8794.html#section-7.3
func Float(b []byte) (float64, error) {
	// A Float Element MUST declare a length of either
	// zero octets (0 bit), four octets (32 bit),
	// or eight octets (64 bit).
	switch len(b) {
	case 0:
		return 0, nil
	case 4:
		return float64(math.Float32frombits(binary.BigEndian.Uint32(b))), nil
	case 8:
		return math.Float64frombits(binary.BigEndian.Uint64(b)), nil
	default:
		return 0, errors.New("ebml: data length must be 0 bit, 32 bit or 64 bit for a float")
	}
}

// String reads a string based on https://www.rfc-editor.org/rfc/rfc8794.html#section-7.5
func String(b []byte) (string, error) {
	return string(b), nil
}

// Date reads a time,Time based on https://www.rfc-editor.org/rfc/rfc8794.html#section-7.6
func Date(b []byte) (time.Time, error) {
	i, err := Int(b)
	if err != nil {
		return time.Time{}, err
	}
	return thirdMillennium.Add(time.Nanosecond * time.Duration(i)), nil
}

type Encoder struct {
	// https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.4
	MaxIDLength uint
	// https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.5
	MaxSizeLength uint
	offset        int64

	buf []byte
	w   io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		MaxIDLength:   8,
		MaxSizeLength: 8,

		buf: make([]byte, 0, newBufferSize),
		w:   w,
	}
}

func (e *Encoder) OutputOffset() int64 {
	return e.offset
}

// WriteElementID writes an Element ID based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-5
//
// TODO: EBMLMaxIDLength can be greater than 8 https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.4
func (e *Encoder) WriteElementID(id schema.ElementID) (n int, err error) {
	if vintAllOne(uint64(id)) {
		return 0, ErrAllOneVINT
	}
	e.buf = e.buf[:e.MaxIDLength]
	w, err := AppendVint(uint64(id), e.buf)
	e.offset += int64(w)
	if err != nil {
		return 0, ErrInvalidVINTWidth
	}
	return e.w.Write(e.buf[:w])
}

// WriteElementDataSize writes an Element Data Size based on
// https://datatracker.ietf.org/doc/html/rfc8794#section-6
//
// TODO: EBMLMaxSizeLength can be greater than 8 https://datatracker.ietf.org/doc/html/rfc8794#section-11.2.5
func (e *Encoder) WriteElementDataSize(ds int64, minW int) (n int, err error) {
	e.buf = e.buf[:e.MaxSizeLength]
	uds := uint64(ds)
	if ds == -1 {
		uds = makeVintDataAllOne(minW)
	}
	w, err := AppendVintData(uds, minW, e.buf)
	e.offset += int64(w)
	if err != nil {
		return 0, ErrInvalidVINTWidth
	}
	return e.w.Write(e.buf[:w])
}

func (e *Encoder) Write(v []byte) (int, error) {
	n, err := e.w.Write(v)
	e.offset += int64(n)
	return n, err
}
