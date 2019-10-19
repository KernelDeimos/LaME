package util

import (
	"strings"
)

type LibUtilString struct {
	Capitalize func(string) string
}

var String LibUtilString

func init() {
	String = LibUtilString{}
	String.Capitalize = strings.Title
}
