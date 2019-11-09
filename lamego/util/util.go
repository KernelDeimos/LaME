package util

import (
	// "strings"

	"github.com/iancoleman/strcase"
)

type LibUtilString struct {
	Capitalize func(string) string
}

var String LibUtilString

func init() {
	String = LibUtilString{}
	String.Capitalize = strcase.ToCamel
}
