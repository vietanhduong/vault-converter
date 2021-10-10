package util

import "reflect"

func IsNullOrEmpty(value interface{}) bool {
	return reflect.ValueOf(value).IsZero()
}