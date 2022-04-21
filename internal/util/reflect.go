package util

import (
	"fmt"
	"reflect"
)

func BuildReflectTypeKey(varType reflect.Type) string {
	const kfmt = "%s.%s.%s"
	if varType == nil {
		return "nil"
	}
	return fmt.Sprintf(
		kfmt,
		varType.PkgPath(),
		varType.Name(),
		varType.String(),
	)
}

func KeyFN(t interface{}) string {
	return BuildReflectTypeKey(reflect.TypeOf(t))
}
