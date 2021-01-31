package ebml

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
	"time"
)

func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

// Decode works like Unmarshal, except it reads the decoder
// stream.
func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer passed to Unmarshal")
	}
	err := d.unmarshal(val.Elem())
	if err == io.EOF {
		err = nil
	}
	return err
}

var (
	timeType        = reflect.TypeOf(time.Time{})
	thirdMillennium = time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
)

// Unmarshal a single EBML element into val.
func (d *Decoder) unmarshal(val reflect.Value) error {
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

	switch v := val; v.Kind() {
	default:
		return errors.New("unknown type " + v.Type().String())

	case reflect.Struct:
		typ := v.Type()
		if typ == timeType {
			t, err := d.readTime()
			if err != nil {
				return err
			}
			v.Set(reflect.ValueOf(t))

			break
		}

		tinfo, err := getTypeInfo(typ)
		if err != nil {
			return err
		}
		for {
			el, err := d.Element()
			if err != nil {
				return err
			}
			dd := NewDecoder(el.Data)
			for i := range tinfo.fields {
				finfo := tinfo.fields[i]
				if el.ID.Val().Cmp(finfo.name) != 0 {
					continue
				}
				f := val.Field(finfo.idx[0])
				err := dd.unmarshal(f)
				if err != nil && err != io.EOF {
					return err
				}
				break
			}
		}

	case reflect.Slice:
		e := v.Type().Elem()
		switch e.Kind() {
		case reflect.Uint8:
			b, err := d.readByteSlice()
			if err != nil {
				return err
			}
			v.SetBytes(b)

		default:
			n := v.Len()
			v.Set(reflect.Append(v, reflect.Zero(e)))
			if err := d.unmarshal(v.Index(n)); err != nil {
				return err
			}
		}

	case reflect.String:
		s, err := d.readString()
		if err != nil {
			return err
		}
		v.SetString(s)

	case reflect.Float32, reflect.Float64:
		f, err := d.readFloat()
		if err != nil {
			return err
		}
		v.SetFloat(f)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := d.readInt()
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := d.readUint()
		if err != nil {
			return err
		}
		v.SetUint(i)
	}

	return nil
}

func (d *Decoder) readByteSlice() ([]byte, error) {
	var err error
	var b []byte
	var bb byte
	for {
		bb, err = d.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		b = append(b, bb)
	}
	return b, err
}

func (d *Decoder) readTime() (time.Time, error) {
	i, err := d.readInt()
	if err != nil {
		return time.Time{}, err
	}
	return thirdMillennium.Add(time.Nanosecond * time.Duration(i)), nil
}

func (d *Decoder) readString() (string, error) {
	b, err := d.readByteSlice()
	if err != nil {
		return "", err
	}
	// TODO: detect value greater than VINTMAX
	return string(b), err
}

func (d *Decoder) readFloat() (float64, error) {
	b, err := d.readByteSlice()
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

func (d *Decoder) readInt() (int64, error) {
	b, err := d.readByteSlice()
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

func (d *Decoder) readUint() (uint64, error) {
	b, err := d.readByteSlice()
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
