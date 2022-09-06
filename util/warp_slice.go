package util

import (
	"reflect"
)

func WrapSlice[T any](slice any) []T {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}
	items := make([]T, v.Len())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		items[i] = item.(T)
	}
	return items
}
