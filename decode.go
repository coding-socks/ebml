package ebml

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"math"
	"reflect"
	"time"
)

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "ebml: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "ebml: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "ebml: Unmarshal(nil " + e.Type.String() + ")"
}

// Decode works like Unmarshal, except it reads the decoder
// stream.
func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	if err := d.decodeRoot(val.Elem(), headerDefinition); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return err
	}
	// TODO: read doctype from header
	bodyDef, err := getDefinition("matroska")
	if err != nil {
		return err
	}
	if err := d.decodeRoot(val.Elem(), bodyDef); err != nil {
		if err == io.EOF {
			err = nil
		}
		return err
	}
	// TODO: decode body / segment
	return nil
}

var (
	timeType        = reflect.TypeOf(time.Time{})
	thirdMillennium = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
)

func findField(val reflect.Value, tinfo *typeInfo, name string) (fieldv reflect.Value, found bool) {
	for i := range tinfo.fields {
		finfo := tinfo.fields[i]
		if name != finfo.name {
			continue
		}
		found = true
		fieldv = val.Field(finfo.idx[0])
		break
	}
	return
}

func (d *Decoder) decodeRoot(val reflect.Value, def Definition) error {
	// Load value from interface, but only if the result will be
	// usefully addressable.
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() {
			val = e
		}
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
	el, err := d.element([]Definition{def})
	if err != nil {
		return err
	}
	if bytes.Compare(el.Definition.ID, CRC32.ID) == 0 {
		return errors.New("ebml: unexpected crc-32 element")
	}
	if bytes.Compare(el.Definition.ID, Void.ID) == 0 {
		// TODO: skip void elements
		return errors.New("ebml: unexpected void element")
	}
	log.Printf("=> 0x%s", el.HexID())
	if val.Kind() != reflect.Struct {
		return errors.New("ebml: unknown root element type: " + val.Type().String())
	}
	typ := val.Type()
	tinfo, err := getTypeInfo(typ)
	if err != nil {
		return err
	}
	fieldv, found := findField(val, tinfo, el.Definition.Name)

	if err := d.decodeSingle(el, found, fieldv, el.Definition.Children); err != nil {
		return err
	}
	log.Printf("<= 0x%s", el.HexID())
	return nil
}

func (d *Decoder) decodeMaster(val reflect.Value, ds dataSize, defs []Definition) error {
	// Load value from interface, but only if the result will be
	// usefully addressable.
	if val.Kind() == reflect.Interface && !val.IsNil() {
		e := val.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() {
			val = e
		}
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}

	start := d.r.Position()
	for {
		// Check size first because it can be 0
		if ds.Known() {
			offset := d.r.Position() - start
			if offset > ds.Size() {
				return errors.New("ebml: element overflow")
			} else if offset == ds.Size() {
				return nil
			}
		}

		el, err := d.element(defs)
		if err != nil {
			if err == errInvalidId {
				continue
			}
			var e *UnknownElementError
			if !ds.Known() && errors.As(err, &e) {
				d.elCache = &e.el
				return nil
			}
			return err
		}
		log.Printf("==> 0x%s", el.HexID())
		switch v := val; v.Kind() {
		default:
			return errors.New("unknown master element type: " + val.Type().String())
		case reflect.Slice:
			// TODO: Consider checking max / min occurrence.
			e := v.Type().Elem()
			n := v.Len()
			v.Set(reflect.Append(v, reflect.Zero(e)))
			val = v.Index(n)
		case reflect.Struct:
			// Everything is ok
		}
		typ := val.Type()
		tinfo, err := getTypeInfo(typ)
		if err != nil {
			return err
		}
		fieldv, found := findField(val, tinfo, el.Definition.Name)

		if err := d.decodeSingle(el, found, fieldv, el.Definition.Children); err != nil {
			return err
		}
		log.Printf("<== 0x%s", el.HexID())
	}
}

func (d *Decoder) decodeSingle(el Element, found bool, val reflect.Value, defs []Definition) error {
	if v := val; v.Kind() == reflect.Slice {
		e := v.Type().Elem()
		// TODO: Consider using `el.Definition.Type != TypeBinary || e.Kind() != reflect.Uint8`.
		if !(el.Definition.Type == TypeBinary && e.Kind() == reflect.Uint8) {
			n := v.Len()
			v.Set(reflect.Append(v, reflect.Zero(e)))
			val = v.Index(n)
		}
	}
	switch el.Definition.Type {
	default:
		return errors.New("ebml: unknown type: " + el.Definition.Type)

	case TypeMaster:
		if err := d.decodeMaster(val, el.DataSize, defs); err != nil {
			return err
		}

	case TypeBinary:
		b, err := d.readByteSlice(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.SetBytes(b)
		}

	case TypeDate:
		t, err := d.readDate(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.Set(reflect.ValueOf(t))
		}

	case TypeFloat:
		f, err := d.readFloat(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.SetFloat(f)
		}

	case TypeInteger:
		i, err := d.readInt(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.SetInt(i)
		}

	case TypeUinteger:
		i, err := d.readUint(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.SetUint(i)
		}

	case TypeString, TypeUTF8:
		str, err := d.readString(el.DataSize.Size())
		if err != nil {
			return err
		}
		if found {
			val.SetString(str)
		}
	}
	return nil
}

func (d *Decoder) readByteSlice(ds int64) ([]byte, error) {
	b := make([]byte, ds)
	_, err := io.ReadFull(d.r, b)
	return b, err
}

func (d *Decoder) readDate(ds int64) (time.Time, error) {
	i, err := d.readInt(ds)
	if err != nil {
		return time.Time{}, err
	}
	return thirdMillennium.Add(time.Nanosecond * time.Duration(i)), nil
}

func (d *Decoder) readString(ds int64) (string, error) {
	b, err := d.readByteSlice(ds)
	if err != nil {
		return "", err
	}
	// TODO: detect value greater than VINTMAX
	return string(b), err
}

func (d *Decoder) readFloat(ds int64) (float64, error) {
	b, err := d.readByteSlice(ds)
	if err != nil {
		return 0, err
	}
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

func (d *Decoder) readInt(ds int64) (int64, error) {
	b, err := d.readByteSlice(ds)
	if err != nil {
		return 0, err
	}
	if len(b) > 8 {
		return 0, errors.New("ebml: max length for an unsigned integer is eight octets")
	}
	i := int64(0)
	for _, bb := range b {
		i = (i << 8) | int64(bb)
	}
	return i, nil
}

func (d *Decoder) readUint(ds int64) (uint64, error) {
	b, err := d.readByteSlice(ds)
	if err != nil {
		return 0, err
	}
	if len(b) > 8 {
		return 0, errors.New("ebml: max length for an unsigned integer is eight octets")
	}
	i := uint64(0)
	for _, bb := range b {
		i = (i << 8) | uint64(bb)
	}
	return i, nil
}
