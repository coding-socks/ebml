package ebml

import (
	"errors"
	"fmt"
	"github.com/coding-socks/ebml/ebmltext"
	"github.com/coding-socks/ebml/schema"
	"io"
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

// ErrElementOverflow signals that an element signals a length
// greater than the parent DataSize.
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
			return nil, fmt.Errorf("ebml: unexpected element %v in root", el.ID)
		case IDVoid:
			if _, err := d.Seek(el.DataSize, io.SeekCurrent); err != nil {
				return nil, fmt.Errorf("ebml: could not skip Void element: %w", err)
			}
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
			return fmt.Errorf("ebml: unexpected element %v in root", el.ID)
		case IDVoid:
			if _, err := d.Seek(el.DataSize, io.SeekCurrent); err != nil {
				return fmt.Errorf("ebml: could not skip Void element: %w", err)
			}
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
	if err == nil && d.elOverflow {
		return ErrElementOverflow
	}
	return err
}

var (
	typeTime      = reflect.TypeOf(time.Time{})
	typeDuration  = reflect.TypeOf(time.Duration(0))
	typeElementID = reflect.TypeOf(schema.ElementID(0))
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
	tinfo, ok := d.typeInfos[typ]
	if !ok {
		var err error
		if tinfo, err = getTypeInfo(typ); err != nil {
			return err
		}
		d.typeInfos[typ] = tinfo
	}

	occurrences := make(map[schema.ElementID]int)
	offsetStart := d.r.InputOffset()
	offset := int64(0)
	for {
		el, _, err := d.NextOf(current, d.r.InputOffset()-offsetStart)
		offset = d.r.InputOffset() - offsetStart
		if err != nil {
			if err == ebmltext.ErrInvalidVINTWidth {
				d.r.Seek(1, io.SeekCurrent)
				continue
			}
			if err == io.EOF {
				break
			}
			var e *UnknownElementError
			if current.DataSize == -1 && errors.As(err, &e) {
				break
			}
			return err
		}
		if current.DataSize != -1 {
			// detect element overflow early to pretend the element is smaller
			if current.DataSize < offset+el.DataSize {
				el.DataSize = current.DataSize - offset
				d.elOverflow = true
			}
			offset += el.DataSize
		}
		def, _ := d.def.Get(el.ID)
		occurrences[el.ID]++
		fieldv, found := findField(val, tinfo, def.Name)
		if !found {
			if el.DataSize != -1 {
				if _, err := d.Seek(el.DataSize, io.SeekCurrent); err != nil {
					return fmt.Errorf("ebml: failed to skip element: %w", err)
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

	if current.DataSize != -1 && offset < current.DataSize {
		return io.ErrUnexpectedEOF
	}

	for sel := range d.def.All() {
		if sel.Default == nil || occurrences[sel.ID] > 0 {
			continue
		}
		fieldv, found := findField(val, tinfo, sel.Name)
		if !found {
			continue
		}
		if fieldv.Kind() == reflect.Ptr {
			if fieldv.IsNil() {
				fieldv.Set(reflect.New(fieldv.Type().Elem()))
			}
			fieldv = fieldv.Elem()
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
			switch v.Type() {
			default:
				return &DecodeTypeError{EBMLType: def.Type, Type: v.Type(), Offset: position}
			case typeElementID:
				// valid type
			}
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

var DefaultAllocationWindow = int64(1<<24) - 1

func (d *Decoder) decodeSingle(el Element, val reflect.Value) error {
	def, _ := d.def.Get(el.ID)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			val.Set(reflect.New(val.Type().Elem()))
		}
		val = val.Elem()
	}
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

	if def.Type == TypeMaster {
		return d.decodeMaster(val, el)
	}

	if int64(cap(d.window)) < el.DataSize {
		n := DefaultAllocationWindow
		for n < el.DataSize {
			n = (n << 1) + 1
		}
		d.window = make([]byte, n)
	}
	b := d.window[:el.DataSize]
	if _, err := io.ReadFull(d.r, b); err != nil {
		return err
	}

	switch def.Type {
	case TypeBinary:
		switch val.Type() {
		default:
			d.window = d.window[el.DataSize:]
			val.SetBytes(b)
		case typeElementID:
			i, err := ebmltext.Uint(b)
			if err != nil {
				return err
			}
			val.SetUint(i)
		}

	case TypeDate:
		t, err := ebmltext.Date(b)
		if err != nil {
			return err
		}
		val.Set(reflect.ValueOf(t))

	case TypeFloat:
		f, err := ebmltext.Float(b)
		if err != nil {
			return err
		}
		val.SetFloat(f)

	case TypeInteger:
		i, err := ebmltext.Int(b)
		if err != nil {
			return err
		}
		val.SetInt(i)

	case TypeUinteger:
		switch val.Type() {
		default:
			i, err := ebmltext.Uint(b)
			if err != nil {
				return err
			}
			val.SetUint(i)
		case typeDuration:
			i, err := ebmltext.Int(b)
			if err != nil {
				return err
			}
			val.SetInt(i)
		}

	case TypeString, TypeUTF8:
		str, err := ebmltext.String(b)
		if err != nil {
			return err
		}
		val.SetString(str)
	}
	return nil
}
