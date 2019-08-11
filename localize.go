/**
 * localize.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package localize

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"regexp"
)

var _ Localizer = &Map{}

// JSVariableRegex matches a valid JavaScript variable name.
// Variable name documentation:
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Grammar_and_types#Variables
var JSVariableRegex = regexp.MustCompile(`^[a-zA-z_\$][a-zA-z_\$0-9]*$`)

// JSReservedRegex matches reserved JavaScript keywords that
// may not be used as variable names. Reserved keyword
// documentation:
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Lexical_grammar#Keywords
var JSReservedRegex = regexp.MustCompile(`^(break|case|catch|class|const|continue|debugger|default|delete|do|else|export|extends|finally|for|function|if|import|in|instanceof|new|return|super|switch|this|throw|try|typeof|var|void|while|with|yield|enum|await|implements|interface|package|private|protected|public|static)$`)

var (
	ErrReservedKeyword     = fmt.Errorf("Reserved variable name provided")
	ErrInvalidVariableName = fmt.Errorf("Invalid variable name provided")
	ErrInvalidKey          = fmt.Errorf("Invalid key name provided")
	ErrInvalidData         = fmt.Errorf("Invalid data provided")

	// ErrNilMap most likely indicates that NewMap() was
	// provided with a nil pointer.
	ErrNilMap = fmt.Errorf("Nil data map field")
)

// Data is an alias for an interface map.
type Data = map[string]interface{}

// Localizer describes a struct that localizes Golang data.
type Localizer interface {
	// Data manipulation.
	Add(key string, data interface{}) error
	Delete(key string) error
	GetData() Data

	// Namespacing.
	SetGlobalName(name string) error
	GetGlobalName() string

	// Localization.
	JS() template.JS
}

// Map takes a set of data, translates it to JavaScript
// primitives, and then formats it for insertion into a global
// browser context.
type Map struct {
	data       Data
	globalName string
}

// NewMap generates a new localization map.
func NewMap(name string, data Data) (*Map, error) {
	if nil == data {
		data = Data{}
	}
	l := &Map{
		data: data,
	}
	if err := l.SetGlobalName(name); nil != err {
		return nil, err
	}
	return l, nil
}

// Add inserts an element with the specified key to the data
// map.
func (l *Map) Add(key string, data interface{}) error {
	if nil == l.data {
		return ErrNilMap
	}
	if "" == key {
		return ErrInvalidKey
	}
	if nil == data {
		return ErrInvalidData
	}

	l.data[key] = data
	if val, ok := l.data[key]; !ok || nil == val {
		return errors.New("Failed to add element")
	}

	return nil
}

// Delete removes an element with the specified key from the
// data map.
func (l *Map) Delete(key string) error {
	if nil == l.data {
		return ErrNilMap
	}
	if "" == key {
		return ErrInvalidKey
	}

	delete(l.data, key)
	if _, ok := l.data[key]; ok {
		return fmt.Errorf(
			"Failed to delete element with key, %v",
			key,
		)
	}

	return nil
}

// GetData retrieves the localization map's data.
func (l *Map) GetData() Data {
	return l.data
}

// SetGlobalName assigns the localization map's global
// JavaScript variable name, which will receive the localized
// data.
func (l *Map) SetGlobalName(name string) error {
	var buf bytes.Buffer
	buf.WriteString(name)
	bytes := buf.Bytes()
	if ok := JSVariableRegex.Match(bytes); !ok {
		return ErrInvalidVariableName
	}
	if ok := JSReservedRegex.Match(bytes); ok {
		return ErrReservedKeyword
	}

	l.globalName = name
	return nil
}

// GetGlobalName retrieves the localization map's global
// JavaScript variable name.
func (l *Map) GetGlobalName() string {
	return l.globalName
}

// JS gets a valid block of template.JS data that represents
// the fields of this Map's "data" field and all its
// children. The returned template.JS block can be directly
// placed into an HTML template (provided by the
// "html/template" package) and output as valid JavaScript
// code.
func (l *Map) JS() template.JS {
	// Generates a buffer that will have the JavaScript
	// string-formatted bytes written to it. The head of the
	// buffer is a global variable assignment.
	buf := bytes.NewBuffer([]byte(fmt.Sprintf("%s = {\n", l.globalName)))

	// Fills the buffer.
	ReflectTarget(reflect.ValueOf(l.data), buf)
	buf.Write([]byte("\n};"))

	return template.JS(buf.String())
}

// ReflectTarget takes a reflect.Value object and recursively
// determines the values of all the fields, sub-fields,
// elements, etc. At each step, the target's type is analyzed
// to see whether or not it's an enclosing type. If the target
// is an enclosing type, then the contents of the target will
// be wrapped appropriately. Square-brackets ("[]") are used
// for translating data to a JavaScript array. Curly-brackets
// ("{}") are used for translating data to a JavaScript object.
// Non-enclosing types simply output according to their
// JavaScript equivalent.
//
// The complete contents of the top-most target is written
// piece-by-piece to the buffer provided.
func ReflectTarget(target reflect.Value, buf *bytes.Buffer) {
	targetType := target.Type().Kind().String()
	switch targetType {
	case "interface":
		f := target.Elem()

		ReflectTarget(f, buf)
	case "struct":
		numFields := target.NumField()
		for i := 0; i < numFields; i++ {
			f := target.Field(i)

			buf.Write([]byte(fmt.Sprintf("\"%s\": {\n", target.Type().Field(i).Name)))
			ReflectTarget(f, buf)
			buf.Write([]byte(fmt.Sprint("},\n")))
		}
	case "map":
		keys := target.MapKeys()
		for _, keyValue := range keys {
			f := target.MapIndex(keyValue)
			fType := f.Type().Kind().String()

			if "map" == fType || "interface" == fType {
				cOpen := "{"
				cClose := "}"

				if "interface" == fType && "map" != f.Elem().Type().Kind().String() {
					cOpen = "["
					cClose = "]"
				}

				buf.Write([]byte(fmt.Sprintf("\"%s\": %s\n", keyValue, cOpen)))
				ReflectTarget(f, buf)
				buf.Write([]byte(fmt.Sprintf("\n%s,\n", cClose)))
			} else {
				buf.Write([]byte(fmt.Sprintf("\"%s\":", keyValue)))
				ReflectTarget(f, buf)
				buf.Write([]byte(fmt.Sprint("\n")))
			}
		}
	case "slice":
		sliceLen := target.Len()
		buf.Write([]byte(fmt.Sprint("[")))
		for i := 0; i < sliceLen; i++ {
			f := target.Index(i)

			ReflectTarget(f, buf)
		}
		buf.Write([]byte(fmt.Sprint("],\n")))
	case "int":
		buf.Write([]byte(fmt.Sprintf("%v,", target.Int())))
	case "string":
		buf.Write([]byte(fmt.Sprintf("\"%v\",", target.String())))
	case "bool":
		buf.Write([]byte(fmt.Sprintf("%v,", target.Bool())))
	case "float64":
		buf.Write([]byte(fmt.Sprintf("%v,", target.Float())))
	}
}
