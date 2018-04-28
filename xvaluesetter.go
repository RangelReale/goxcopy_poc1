package xcopy

import (
	"errors"
	"reflect"
)

type XValueSetter interface {
	//GetValueType() reflect.Type
	SetValue(v reflect.Value) error
	//HasFields() bool
	//SetField(fieldname string, v reflect.Value) error
}

//
// Error
//

type XValueSetter_Error struct {
}

func (s *XValueSetter_Error) GetValueType() reflect.Type {
	return reflect.TypeOf(nil)
}

func (s *XValueSetter_Error) SetValue(v reflect.Value) error {
	return errors.New("Value cannot be set")
}

func (s *XValueSetter_Error) HasFields() bool {
	return false
}

func (s *XValueSetter_Error) SetField(fieldname string, v reflect.Value) error {
	return errors.New("Field value cannot be set")
}

//
// StructField
//

type XValueSetter_StructField struct {
	v         reflect.Value
	fieldname string
}

func (s *XValueSetter_StructField) SetValue(v reflect.Value) error {
	xv := reflect.Indirect(s.v)

	fieldValue := xv.FieldByName(s.fieldname)
	fieldType, fieldFound := xv.Type().FieldByName(s.fieldname)

	if !fieldFound {
		return nil
	}

	// if destination struct is a nil pointer, create an instance
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldType.Type.Elem()))
	}

	fieldValue.Set(v)
	return nil
}

//
// Slice
//

type XValueSetter_Slice struct {
	v   reflect.Value
	idx int
}

func (s *XValueSetter_Slice) SetValue(v reflect.Value) error {
	xv := reflect.Indirect(s.v)

	for s.idx >= xv.Len() {
		// add zero values until the index
		if xv.Len() == 0 {
			xv.Set(reflect.MakeSlice(reflect.SliceOf(xv.Type().Elem()), 0, 0))
		}
		xv.Set(reflect.Append(xv.Slice(0, xv.Len()), reflect.Zero(xv.Type().Elem())))
	}

	xv.Index(s.idx).Set(v)
	return nil
}

//
// Map
//

type XValueSetter_Map struct {
	v   reflect.Value
	key reflect.Value
}

func (s *XValueSetter_Map) SetValue(v reflect.Value) error {
	reflect.Indirect(s.v).SetMapIndex(s.key, v)
	return nil
}
