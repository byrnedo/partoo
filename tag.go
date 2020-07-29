package partoo

import (
	"reflect"
)

// dumbStructToHeader converts a struct of string fields (with tag 'header') to a http.Header map
type NamedField struct {
	Name string
	Field interface{}
}

type NamedFields []NamedField

func (cm NamedFields) Names() (ret ColNames) {
	idx := 0
	ret = make([]string, len(cm))
	for _, v := range cm {
		ret[idx] = v.Name
		idx++
	}
	return
}


func (cm NamedFields) Fields() (ret []interface{}) {
	idx := 0
	ret = make([]interface{}, len(cm))
	for _, v := range cm {
		ret[idx] = v.Field
		idx++
	}
	return
}

func getColumnNames(table Table) (ret []NamedField) {
	t := reflect.TypeOf(table)

	v := reflect.ValueOf(table)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for _, col := range table.Columns() {

		sf := findStructField(v, reflect.ValueOf(col))

		colName := sf.Tag.Get("sql")
		if colName == "" {
			panic("struct field must have `sql` tag if included in Columns() output")
		}
		ret = append(ret, NamedField{Name: colName, Field: col})
	}
	return
}

// Author: Copied from https://github.com/go-ozzo/ozzo-validation/
//
// findStructField looks for a field in the given struct.
// The field being looked for should be a pointer to the actual struct field.
// If found, the field info will be returned. Otherwise, nil will be returned.
func findStructField(structValue reflect.Value, fieldValue reflect.Value) *reflect.StructField {
	ptr := fieldValue.Pointer()
	for i := structValue.NumField() - 1; i >= 0; i-- {
		sf := structValue.Type().Field(i)
		if ptr == structValue.Field(i).UnsafeAddr() {
			// do additional type comparison because it's possible that the address of
			// an embedded struct is the same as the first field of the embedded struct
			if sf.Type == fieldValue.Elem().Type() {
				return &sf
			}
		}
		if sf.Anonymous {
			// delve into anonymous struct to look for the field
			fi := structValue.Field(i)
			if sf.Type.Kind() == reflect.Ptr {
				fi = fi.Elem()
			}
			if fi.Kind() == reflect.Struct {
				if f := findStructField(fi, fieldValue); f != nil {
					return f
				}
			}
		}
	}
	return nil
}