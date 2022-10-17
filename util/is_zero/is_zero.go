package is_zero

import "reflect"

func Check(any interface{}) bool {
	return reflect.DeepEqual(any, reflect.Zero(reflect.TypeOf(any)).Interface())
}

func CheckComparable[T comparable](any T) bool {
	var t T
	return t == any
}
