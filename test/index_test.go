/**
 * index_test.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package test

import (
	"html/template"
	"testing"

	"github.com/foresthoffman/localize"
)

type testCase struct {
	Input    localize.Data
	Expected []template.JS
}

var testCases = map[string]testCase{
	// Int case.
	"intCase": testCase{
		Input: localize.Data{
			"int": 1954,
		},
		Expected: []template.JS{template.JS(
			`intCase = {
"int": [
1954,
],

};`,
		)},
	},
	// Int array case.
	"intArrayCase": testCase{
		Input: localize.Data{
			"intArray": []int{1, 2, 3, 4, 5},
		},
		Expected: []template.JS{template.JS(
			`intArrayCase = {
"intArray": [
[1,2,3,4,5,],

],

};`,
		)},
	},
	// Multi-dimensional array case.
	"multiArrayCase": testCase{
		Input: localize.Data{
			"arrayArray": [][]int{
				[]int{6, 7, 8, 9, 10},
				[]int{11, 12, 13, 14, 15},
			},
		},
		Expected: []template.JS{template.JS(
			`multiArrayCase = {
"arrayArray": [
[[6,7,8,9,10,],
[11,12,13,14,15,],
],

],

};`,
		)},
	},
	// Map case.
	"mapCase": testCase{
		Input: localize.Data{
			"assocArray": map[string]string{
				"baz": "fubar",
				"foo": "bar",
			},
		},
		Expected: []template.JS{
			template.JS(
				`mapCase = {
"assocArray": {
"baz":"fubar",
"foo":"bar",

},

};`,
			),
			template.JS(
				`mapCase = {
"assocArray": {
"foo":"bar",
"baz":"fubar",

},

};`,
			),
		},
	},
}
var maps = make(map[string]*localize.Map)

// TestNewMap ensures that the data maps can be instantiated as
// expected.
func TestNewMap(t *testing.T) {
	for name, tCase := range testCases {
		m, err := localize.NewMap(name, tCase.Input)
		if nil != err {
			t.Run(name, func(t *testing.T) {
				t.Errorf("Failed to create new map,\nerr: %v\n", err)
			})
		}
		maps[name] = m
	}
}

// TestJS ensures that the localized data is formatted as
// expected.
func TestJS(t *testing.T) {
	for name, m := range maps {
		output := m.JS()

		matched := false
		for _, expected := range testCases[name].Expected {
			if expected == output {
				matched = true
				break
			}
		}
		if !matched {
			t.Run(name, func(t *testing.T) {
				t.Errorf("Expected one of: %v,\ngot: %q\n", testCases[name].Expected, output)
			})
		}
	}
}

// TestInvalidVariableName ensures that invalid variable names
// cannot be assigned to the localized data.
func TestInvalidVariableName(t *testing.T) {
	invalidCases := map[string]testCase{
		"2var": testCase{Input: localize.Data{}},
		"-var": testCase{Input: localize.Data{}},
		"*var": testCase{Input: localize.Data{}},
	}
	for name, tCase := range invalidCases {
		if _, err := localize.NewMap(name, tCase.Input); localize.ErrInvalidVariableName != err {
			t.Run(name, func(t *testing.T) {
				t.Errorf("Expected err: %v,\ngot: %v\n", localize.ErrInvalidVariableName, err)
			})
		}
	}
}

// TestReservedVariableName ensures that reserved variable
// names cannot be assigned to the localized data.
func TestReservedVariableName(t *testing.T) {
	reservedCases := map[string]testCase{
		"var":      testCase{Input: localize.Data{}},
		"function": testCase{Input: localize.Data{}},
		"await":    testCase{Input: localize.Data{}},
		"import":   testCase{Input: localize.Data{}},
	}
	for name, tCase := range reservedCases {
		if _, err := localize.NewMap(name, tCase.Input); localize.ErrReservedKeyword != err {
			t.Run(name, func(t *testing.T) {
				t.Errorf("Expected err: %v,\ngot: %v\n", localize.ErrReservedKeyword, err)
			})
		}
	}
}
