package lamelib

import (
	"testing"

	"github.com/KernelDeimos/LaME/lamego/lamelib/l"
)

func TestIndexOf(t *testing.T) {
	strlib := l.String{}
	i := strlib.IndexOf("fuzzywuzzywasabear", "wuzzy")
	if i != 5 {
		t.Error("expected index 5")
	}
	i = strlib.IndexOf("awuzzy", "wuzzy")
	if i != 1 {
		t.Error("expected index 1")
	}
	i = strlib.IndexOf("wuzzy", "wuzzy")
	if i != 0 {
		t.Error("expected index 0")
	}
	i = strlib.IndexOf("wuzz", "wuzzy")
	if i != -1 {
		t.Error("expected index -1")
	}
}
