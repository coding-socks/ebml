package ebml

import (
	"errors"
	"math/big"
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
	name    *big.Int
	parents []string
}

// getTypeInfo returns the typeInfo structure with details necessary
// for marshaling and unmarshaling typ.
func getTypeInfo(typ reflect.Type) (*typeInfo, error) {
	// TODO: use cache to load typeInfo

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
						if err := addFieldInfo(typ, tinfo, &finfo); err != nil {
							return nil, err
						}
					}
					continue
				}
			}

			finfo, err := structFieldInfo(typ, &f)
			if err != nil {
				return nil, err
			}

			// Add the field if it doesn't conflict with other fields.
			if err := addFieldInfo(typ, tinfo, finfo); err != nil {
				return nil, err
			}
		}
	}
	return tinfo, nil
}

// structFieldInfo builds and returns a fieldInfo for f.
func structFieldInfo(typ reflect.Type, f *reflect.StructField) (*fieldInfo, error) {
	finfo := &fieldInfo{idx: f.Index}

	tag := f.Tag.Get("ebml")

	// TODO: Parse flags here. No flags for now.

	// TODO: Parse nested structure. Not allowed for now.
	name, ok := big.NewInt(0).SetString(tag, 0)
	if !ok {
		return nil, errors.New("ebml: cannot read ID from tag")
	}
	finfo.name = name

	return finfo, nil
}

// addFieldInfo adds finfo to tinfo.fields if there are no
// conflicts, or if conflicts arise from previous fields that were
// obtained from deeper embedded structures than finfo. In the latter
// case, the conflicting entries are dropped.
// A conflict occurs when the path (parent + name) to a field is
// itself a prefix of another path, or when two paths match exactly.
// It is okay for field paths to share a common, shorter prefix.
func addFieldInfo(typ reflect.Type, tinfo *typeInfo, newf *fieldInfo) error {
	var conflicts []int
	// TODO: figure all conflicts

	// Without conflicts, add the new field and return.
	if conflicts == nil {
		tinfo.fields = append(tinfo.fields, *newf)
		return nil
	}

	// TODO: handle conflicts

	// TODO: implement addFieldInfo
	return errors.New("ebml: conflict handling is not implemented")
}
