package util

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

func IsNullOrEmpty(value interface{}) bool {
	return reflect.ValueOf(value).IsZero()
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
