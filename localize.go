/**
 * localize.go
 *
 * Copyright (c) 2017 Forest Hoffman. All Rights Reserved.
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
)

type LocalizeMap struct {
	data    map[string]interface{}
	varName string
}

// Generates a localization map.
//
// The data parameter is a map of interfaces, which will be type-asserted to types that will
// correspond to valid JavaScript primitives.
//
// The LocalizeMap type has a "varName" field which contains the name of the global JavaScript
// variable that will contain all the localized data. See "localize.SetVarName()" for setting this
// variable name. By default, the variable name will be "_globalVars".
func NewMap(data *map[string]interface{}) *LocalizeMap {
	if nil == data {
		data = &map[string]interface{}{}
	}
	return &LocalizeMap{
		data:    *data,
		varName: "_globalVars",
	}
}

// Sets the localization map's JavaScript variable name, which will receive the localized data.
// The object is returned upon success, with a nil error. Error will be non-nil when the provided
// name is an empty string.
func (l *LocalizeMap) SetVarName(name string) (*LocalizeMap, error) {
	if "" == name {
		return nil, errors.New("LocalizeMap.SetVarName(): The name cannot be an empty string.")
	}
	l.varName = name
	return l, nil
}

// Adds an element with a specified key to the data map. The error is nil upon success.
func (l *LocalizeMap) Add(key string, data interface{}) error {
	if nil == l.data {
		return errors.New("LocalizeMap.Add(): Cannot add element to a nil map.")
	}

	if "" == key {
		return errors.New("LocalizeMap.Add(): Cannot add value to an empty key.")
	}

	if nil == data {
		return errors.New("LocalizeMap.Add(): Cannot add element with nil data.")
	}

	l.data[key] = data
	if val, ok := l.data[key]; !ok || nil == val {
		return errors.New("LocalizeMap.Add(): Failed to add element.")
	}

	return nil
}

// Deletes an element with a specified key from the data map. The error is nil upon success.
func (l *LocalizeMap) Delete(key string) error {
	if nil == l.data {
		return errors.New("LocalizeMap.Delete(): Cannot delete element from a nil map.")
	}

	if "" == key {
		return errors.New("LocalizeMap.Delete(): Cannot delete element with an empty key.")
	}

	delete(l.data, key)
	if _, ok := l.data[key]; ok {
		return errors.New(
			fmt.Sprintf(
				"LocalizeMap.Delete(): Failed to delete element with key, %v.",
				key,
			),
		)
	}

	return nil
}

// Retrieves the localization map's data.
func (l *LocalizeMap) GetData() map[string]interface{} {
	return l.data
}

// Retrieves the localization map's JavaScript variable name, which will receive the localized data.
func (l *LocalizeMap) GetVarName() string {
	return l.varName
}

// Gets a valid block of template.JS data that represents the fields of this LocalizeMap's "data"
// field and all its children. The returned template.JS block can be directly placed into an
// HTML template (provided by the "html/template" package) and output as valid JavaScript code.
func (l *LocalizeMap) JS() template.JS {

	// generates a buffer that will have the JS bytes written to it, and places a global variable
	// at the beginning for the data to be assigned to
	buf := bytes.NewBuffer([]byte(fmt.Sprintf("%s = {\n", l.varName)))

	// fills the buffer
	ReflectTarget(reflect.ValueOf(l.data), buf)
	buf.Write([]byte("\n};"))

	return template.JS(buf.String())
}

// Takes a reflect.Value object and recursively determines the values of all the fields, sub-fields,
// elements, etc. At each step, the target's type is analyzed to see whether or not it's an enclosing
// type. If the target is an enclosing type, then the contents of the target will be wrapped
// appropriately. Square-brackets ("[]") are used for translating data to a JavaScript
// array. Curly-brackets ("{}") are used for translating data to a JavaScript object. Non-enclosing
// types simply output according to their JavaScript equivalent.
//
// The complete contents of the top-most target is written piece-by-piece to the buffer provided.
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
