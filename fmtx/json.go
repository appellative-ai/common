package fmtx

import "fmt"

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

// JsonString - Json format a string value
func JsonString(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "\"" + value + "\""
}

// JsonMarkup - markup a name/value pair
func JsonMarkup(name, value string, stringValue bool) string {
	if len(value) == 0 {
		return fmt.Sprintf(markupNull, name)
	}
	format := markupString
	if !stringValue {
		format = markupValue
	}
	return fmt.Sprintf(format, name, value)
}
