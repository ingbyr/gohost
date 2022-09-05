package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"reflect"
)

type itemWrapper func(any) list.Item

func wrapListItems(slice any) []list.Item {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}
	items := make([]list.Item, v.Len())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		items[i] = item.(list.Item)
	}
	return items
}
