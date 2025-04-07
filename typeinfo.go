package ebml

import (
	"reflect"
)

// typeInfo holds details for the ebml representation of a type.
type typeInfo struct {
	ebmlID *fieldInfo
	fields []fieldInfo
}

// fieldInfo holds details for the ebml representation of a single field.
type fieldInfo struct {
	idx     []int
	name    string
	parents []string
}

// getTypeInfo returns the typeInfo structure with details necessary
// for marshaling and unmarshaling typ.
func getTypeInfo(typ reflect.Type) (*typeInfo, error) {
	tinfo := &typeInfo{}
	if typ.Kind() == reflect.Struct {
		n := typ.NumField()
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if (f.PkgPath != "" && !f.Anonymous) || f.Tag.Get("ebml") == "-" {
				continue // Private field
			}

			if f.Anonymous {
				t := f.Type
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				if t.Kind() == reflect.Struct {
					inner, err := getTypeInfo(t)
					if err != nil {
						return nil, err
					}
					if tinfo.ebmlID == nil {
						tinfo.ebmlID = inner.ebmlID
					}
					for _, finfo := range inner.fields {
						finfo.idx = append([]int{i}, finfo.idx...)
						tinfo.fields = append(tinfo.fields, finfo)
					}
					continue
				}
			}

			finfo, err := structFieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}

			// Add the field if it doesn't conflict with other fields.
			tinfo.fields = append(tinfo.fields, *finfo)
		}
	}
	return tinfo, nil
}

// structFieldInfo builds and returns a fieldInfo for f.
func structFieldInfo(typ reflect.Type, f *reflect.StructField) (*fieldInfo, error) {
	finfo := &fieldInfo{idx: f.Index}

	tag := f.Tag.Get("ebml")

	name := tag
	if name == "" {
		name = f.Name
	}
	finfo.name = name

	return finfo, nil
}

var (
	TypeInteger  = "integer"
	TypeUinteger = "uinteger"
	TypeFloat    = "float"
	TypeString   = "string"
	TypeDate     = "date"
	TypeUTF8     = "utf-8"
	TypeMaster   = "master"
	TypeBinary   = "binary"
)
