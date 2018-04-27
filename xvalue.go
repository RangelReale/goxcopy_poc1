package xcopy

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/RangelReale/rprim"
)

// XValue interface
type XValue interface {
	Name() string
	IsXValue()
	To(dst interface{}) error
	ToXValue(dst XValue) error

	HasFields() bool
	SetField(fieldname string, v reflect.Value) error
}

//
// Struct
//
type XValue_Struct struct {
	v reflect.Value
}

func (x *XValue_Struct) Name() string {
	return "Struct"
}

func (x *XValue_Struct) IsXValue() {}

func (x *XValue_Struct) To(dst interface{}) error {
	return x.ToXValue(XValueOfInterface(dst))
}

func (x *XValue_Struct) ToXValue(dst XValue) error {
	if !dst.HasFields() {
		return fmt.Errorf("Cannot copy Struct to %s", dst.Name())
	}

	xv := reflect.Indirect(x.v)

	for i := 0; i < xv.NumField(); i++ {
		var (
			vField = xv.Field(i)
			tField = xv.Type().Field(i)
		)

		// Is exportable?
		if tField.PkgPath != "" {
			continue
		}

		err := dst.SetField(tField.Name, vField)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *XValue_Struct) HasFields() bool {
	return true
}

func (x *XValue_Struct) SetField(fieldname string, v reflect.Value) error {
	xv := reflect.Indirect(x.v)

	var (
		fieldValue            = xv.FieldByName(fieldname)
		fieldType, fieldFound = xv.Type().FieldByName(fieldname)
	)

	if !fieldFound {
		return nil
	}

	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldType.Type.Elem()))
	}

	return XValueOfValue(v).ToXValue(XValueOfValue(fieldValue))
}

//
// Slice
//
type XValue_Slice struct {
	v reflect.Value
}

func (x *XValue_Slice) Name() string {
	return "Slice"
}

func (x *XValue_Slice) IsXValue() {}

func (x *XValue_Slice) To(dst interface{}) error {
	return x.ToXValue(XValueOfInterface(dst))
}

func (x *XValue_Slice) ToXValue(dst XValue) error {
	if !dst.HasFields() {
		return fmt.Errorf("Cannot copy Slice to %s", dst.Name())
	}

	xv := reflect.Indirect(x.v)

	for i := 0; i < xv.Len(); i++ {
		var (
			vField = xv.Index(i)
		)

		err := dst.SetField(strconv.Itoa(i), vField)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *XValue_Slice) HasFields() bool {
	return true
}

func (x *XValue_Slice) SetField(fieldname string, v reflect.Value) error {
	idx, err := strconv.ParseInt(fieldname, 10, 32)
	if err != nil {
		return fmt.Errorf("Error parsing slice index '%s': %v", fieldname, err)
	}

	xv := reflect.Indirect(x.v)

	for int(idx) >= xv.Len() {
		// add zero values until the index
		if xv.Len() == 0 {
			xv.Set(reflect.MakeSlice(reflect.SliceOf(xv.Type().Elem()), 0, 0))
		}
		xv.Set(reflect.Append(xv.Slice(0, xv.Len()), reflect.Zero(xv.Type().Elem())))
	}

	var (
		fieldValue = xv.Index(int(idx))
	)

	return XValueOfValue(v).ToXValue(XValueOfValue(fieldValue))
}

//
// Primitive
//
type XValue_Primitive struct {
	v reflect.Value
}

func (x *XValue_Primitive) Name() string {
	return "Primitive"
}

func (x *XValue_Primitive) IsXValue() {}

func (x *XValue_Primitive) To(dst interface{}) error {
	return x.ToXValue(XValueOfInterface(dst))
}

func (x *XValue_Primitive) ToXValue(dst XValue) error {
	if dst.HasFields() {
		return fmt.Errorf("Cannot copy Primitive to %s", dst.Name())
	}

	switch xdst := dst.(type) {
	case *XValue_Primitive:
		return x.setPrimitiveValue(xdst)
	default:
		return fmt.Errorf("Cannot copy Primitive to %s", dst.Name())
	}

	return nil
}

func (x *XValue_Primitive) HasFields() bool {
	return false
}

func (x *XValue_Primitive) SetField(fieldname string, v reflect.Value) error {
	return fmt.Errorf("Cannot set Field on Primitive")
}

func (x *XValue_Primitive) setPrimitiveValue(dst *XValue_Primitive) error {
	cop := rprim.ConvertOp(dst.v, x.v, rprim.COP_ALLOW_NIL_TO_ZERO)
	if cop == nil {
		return fmt.Errorf("Could not convert between primitives %s and %s", x.v.Kind().String(), dst.v.Kind().String())
	}
	cv, err := cop(x.v, dst.v.Type())
	if err != nil {
		return err
	}
	dst.v.Set(cv)
	return nil
}
