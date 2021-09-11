package ebml

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/schema"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

// An DecodeTypeError describes an EBML value that was
// not appropriate for a value of a specific Go type.
type DecodeTypeError struct {
	EBMLType string       // description of EBML type - "integer", "binary", "master"
	Type     reflect.Type // type of Go value it could not be assigned to
	Offset   int64        // error occurred after reading Offset bytes
	Path     string       // the full path from root node to the field
}

func (e *DecodeTypeError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("ebml: cannot unmarshal %s into Go struct field %s of type %s", e.EBMLType, e.Path, e.Type)
	}
	return fmt.Sprintf("ebml: cannot unmarshal %s into Go value of type %s", e.EBMLType, e.Type)
}

func (e *DecodeTypeError) extendError(p string) {
	if e.Path == "" {
		e.Path = p
		return
	}
	e.Path = p + "." + e.Path
}

// An InvalidDecodeError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidDecodeError struct {
	Type reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.Type == nil {
		return "ebml: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "ebml: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "ebml: Unmarshal(nil " + e.Type.String() + ")"
}

var ErrElementOverflow = errors.New("ebml: element overflow")

// DecodeHeader decodes the document header.
func (d *Decoder) DecodeHeader() (*EBML, error) {
	for {
		el, _, err := d.Next()
		if err != nil {
			return nil, err
		}
		switch el.ID {
		default:
			return nil, fmt.Errorf("ebml: unexpected element %s in root", el.ID)
		case IDVoid:
			continue
		case IDEBML:
			d.def = HeaderDef
			d.r.MaxIDLength = DefaultMaxIDLength
			d.r.MaxSizeLength = DefaultMaxSizeLength
			var h EBML
			err := d.Decode(&h)
			if err != nil {
				return nil, err
			}
			d.def, err = Definition(h.DocType)
			if err != nil {
				return nil, err
			}
			d.r.MaxIDLength = h.EBMLMaxIDLength
			d.r.MaxSizeLength = h.EBMLMaxSizeLength
			return &h, err
		}
	}
}

// DecodeBody decodes the EBML Body and stores the result in the value
// pointed to by v. If v is nil or not a pointer, DecodeBody returns
// an InvalidDecodeError.
func (d *Decoder) DecodeBody(v interface{}) error {
	for {
		el, _, err := d.Next()
		if err != nil {
			return err
		}
		switch el.ID {
		default:
			return fmt.Errorf("ebml: unexpected element %s in root", el.ID)
		case IDVoid:
			continue
		case d.def.Root.ID:
			return d.Decode(v)
		}
	}
}

func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return &InvalidDecodeError{reflect.TypeOf(v)}
	}
	if d.el == nil {
		return fmt.Errorf("ebml: missing decoded element (forgotten call Next?)")
	}
	err := d.decodeSingle(*d.el, val.Elem())
	d.el = nil
	return err
}

var (
	typeTime     = reflect.TypeOf(time.Time{})
	typeDuration = reflect.TypeOf(time.Duration(0))

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

func (d *Decoder) decodeMaster(val reflect.Value, current Element) error {
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

	occurrences := make(map[string]int)
	var offset int64
	for {
		// Check end of element first because data size can be 0
		if end, err := d.EndOfElement(current, offset); err != nil {
			return err
		} else if end {
			break
		}

		el, n, err := d.r.Next()
		if current.DataSize.Known() {
			offset += int64(n)
		}
		if err != nil {
			if err == ErrInvalidVINTLength {
				continue
			}
			var e *UnknownElementError
			if !current.DataSize.Known() && errors.As(err, &e) {
				break
			}
			return err
		}
		if current.DataSize.Known() {
			offset += el.DataSize.Size()
		}
		def, _ := d.def.Get(el.ID)
		occurrences[el.ID]++
		fieldv, found := findField(val, tinfo, def.Name)
		if !found {
			if el.DataSize.Known() {
				if _, err := d.Seek(el.DataSize.Size(), io.SeekCurrent); err != nil {
					return fmt.Errorf("ebml: was not able to skip element: %w", err)
				}
				continue
			} else if def.Type == TypeMaster {
				if err := d.decodeMaster(val, el); err != nil {
					return err
				}
				continue
			} else {
				return errors.New("ebml: only a master element is allowed to be of unknown size")
			}
		}

		if err := d.decodeSingle(el, fieldv); err != nil {
			if e, ok := err.(*DecodeTypeError); ok {
				e.extendError(val.Type().Name())
			}
			return err
		}
	}

	elements := d.def.Values()
	for i := range elements {
		sel := elements[i]
		if sel.Default == nil || occurrences[sel.ID] > 0 {
			continue
		}
		if sel.Type == TypeMaster {
			// TODO: catch this when Doc Type is registered.
			panic("ebml: master Elements MUST NOT declare a default value.")
		}
		fieldv, found := findField(val, tinfo, sel.Name)
		if !found {
			continue
		}
		if v := fieldv; v.Kind() == reflect.Slice {
			e := v.Type().Elem()
			if !(sel.Type == TypeBinary && e.Kind() == reflect.Uint8) {
				n := v.Len()
				v.Set(reflect.Append(v, reflect.Zero(e)))
				fieldv = v.Index(n)
			}
		}

		if err := validateReflectType(fieldv, sel, 0); err != nil {
			if e, ok := err.(*DecodeTypeError); ok {
				e.extendError(sel.Name)
				e.extendError(val.Type().Name())
			}
			return err
		}
		switch sel.Type {
		case TypeInteger:
			x, _ := strconv.ParseInt(*sel.Default, 10, 64)
			fieldv.SetInt(x)
		case TypeUinteger:
			if fieldv.Type() == typeDuration {
				x, _ := strconv.ParseInt(*sel.Default, 10, 64)
				fieldv.SetInt(x)
			} else {
				x, _ := strconv.ParseUint(*sel.Default, 10, 64)
				fieldv.SetUint(x)
			}
		case TypeFloat:
			x, _ := strconv.ParseFloat(*sel.Default, 64)
			fieldv.SetFloat(x)
		case TypeString:
			fieldv.SetString(*sel.Default)
		default:
			return fmt.Errorf("default not supported: %s", sel.Type)
		}
	}
	return nil
}

func validateReflectType(v reflect.Value, def schema.Element, position int64) error {
	switch def.Type {
	default:
		return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}

	case TypeMaster:
		switch v.Kind() {
		default:
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		case reflect.Struct:
			// valid type
		}

	case TypeBinary:
		switch v.Kind() {
		default:
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		case reflect.Slice:
			e := v.Type().Elem()
			if e.Kind() != reflect.Uint8 {
				return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
			}
		}

	case TypeDate:
		switch v.Type() {
		default:
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		case typeTime:
			// valid type
		}

	case TypeFloat:
		switch v.Kind() {
		default:
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		case reflect.Float32, reflect.Float64:
			// valid type
		}

	case TypeInteger:
		switch v.Kind() {
		default:
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		case reflect.Int, reflect.Int64, reflect.Int32:
			// valid type
		}

	case TypeUinteger:
		switch v.Kind() {
		default:
			switch v.Type() {
			default:
				return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
			case typeDuration:
				// valid type
			}
		case reflect.Uint, reflect.Uint64, reflect.Uint32:
			// valid type
		}

	case TypeString, TypeUTF8:
		if v.Kind() != reflect.String {
			return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
		}
	}
	return nil
}

func (d *Decoder) decodeSingle(el Element, val reflect.Value) error {
	def, _ := d.def.Get(el.ID)
	if v := val; v.Kind() == reflect.Slice {
		e := v.Type().Elem()
		if !(def.Type == TypeBinary && e.Kind() == reflect.Uint8) {
			n := v.Len()
			v.Set(reflect.Append(v, reflect.Zero(e)))
			val = v.Index(n)
		}
	}
	if err := validateReflectType(val, def, 0); err != nil {
		if e, ok := err.(*DecodeTypeError); ok {
			e.extendError(def.Name)
		}
		return err
	}

	switch def.Type {
	case TypeMaster:
		if err := d.decodeMaster(val, el); err != nil {
			return err
		}

	case TypeBinary:
		b, err := d.readByteSlice(el.DataSize.Size())
		if err != nil {
			return err
		}
		val.SetBytes(b)

	case TypeDate:
		t, err := d.readDate(el.DataSize.Size())
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(t))

	case TypeFloat:
		f, err := d.readFloat(el.DataSize.Size())
		if err != nil {
			return err
		}
		val.SetFloat(f)

	case TypeInteger:
		i, err := d.readInt(el.DataSize.Size())
		if err != nil {
			return err
		}
		val.SetInt(i)

	case TypeUinteger:
		switch val.Type() {
		default:
			i, err := d.readUint(el.DataSize.Size())
			if err != nil {
				return err
			}
			val.SetUint(i)
		case typeDuration:
			i, err := d.readInt(el.DataSize.Size())
			if err != nil {
				return err
			}
			val.SetInt(i)
		}

	case TypeString, TypeUTF8:
		str, err := d.readString(el.DataSize.Size())
		if err != nil {
			return err
		}
		val.SetString(str)
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
