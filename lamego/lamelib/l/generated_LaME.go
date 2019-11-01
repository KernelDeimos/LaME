// GENERATED CODE - changes to this file may be overwritten

package l

import "encoding/json"

type String struct {
}
func (o *String) indexOf(subject string,substr string) int {
	var i int
	var e int
	var testval string
	var lensubject int
	var lensubstr int
	
	lensubject = len(subject)
	lensubstr = len(substr)
	if lensubject == 0 {
		return 0
	}
	if lensubject < lensubstr {
		return -1
	}
	i = 0
	e = lensubject - lensubstr
	for i < e {
		testval = (subject)[(i):(i + lensubstr)]
		if testval == substr {
			return i
		}
		i = i + 1
	}
	return -1
}
func (o *String) serializeJSON() string {
	
	return (func() string {
		bout, err := json.Marshal(o)
		if err != nil { return "" }
		return string(bout)
	})()
}
