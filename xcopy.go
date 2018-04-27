package xcopy

import (
	"fmt"
	"reflect"
)

// Functions
func XValueOfInterface(v interface{}) XValue {
	return XValueOfValue(reflect.ValueOf(v))
}

func XValueOfValue(rv reflect.Value) XValue {
	riv := rv.Type()
	if rv.Kind() == reflect.Ptr {
		riv = riv.Elem()
	}

	switch riv.Kind() {
	case reflect.Struct:
		return &XValue_Struct{rv}
	case reflect.String, reflect.Int:
		return &XValue_Primitive{rv}
	case reflect.Slice:
		return &XValue_Slice{rv}
	default:
		panic(fmt.Sprintf("Unknown type %s", riv.Kind().String()))
	}
	return nil
}

// Main function
func XCopy(src interface{}) XValue {
	return XValueOfInterface(src)
}
