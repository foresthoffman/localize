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
	data     localize.Data
	expected []template.JS
}

var testCases = map[string]testCase{
	// Int case.
	"intCase": testCase{
		data: localize.Data{
			"int": 1954,
		},
		expected: []template.JS{template.JS(
			`intCase = {
"int": [
1954,
],

};`,
		)},
	},
	// Int array case.
	"intArrayCase": testCase{
		data: localize.Data{
			"intArray": []int{1, 2, 3, 4, 5},
		},
		expected: []template.JS{template.JS(
			`intArrayCase = {
"intArray": [
[1,2,3,4,5,],

],

};`,
		)},
	},
	// Multi-dimensional array case.
	"multiArrayCase": testCase{
		data: localize.Data{
			"arrayArray": [][]int{
				[]int{6, 7, 8, 9, 10},
				[]int{11, 12, 13, 14, 15},
			},
		},
		expected: []template.JS{template.JS(
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
		data: localize.Data{
			"assocArray": map[string]string{
				"baz": "fubar",
				"foo": "bar",
			},
		},
		expected: []template.JS{
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
		m, err := localize.NewMap(name, tCase.data)
		if nil != err {
			t.Fatalf("Failed to create new map for case, %q,\nerr: %v\n", name, err)
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
		for _, expected := range testCases[name].expected {
			if expected == output {
				matched = true
				break
			}
		}
		if !matched {
			t.Fatalf("Expected one of: %v,\ngot: %q\n", testCases[name].expected, output)
		}
	}
}

// TestInvalidVariableName ensures that invalid variable names
// cannot be assigned to the localized data.
func TestInvalidVariableName(t *testing.T) {
	invalidCases := map[string]testCase{
		"2var": testCase{data: localize.Data{}},
		"-var": testCase{data: localize.Data{}},
		"*var": testCase{data: localize.Data{}},
	}
	for name, tCase := range invalidCases {
		if _, err := localize.NewMap(name, tCase.data); localize.ErrInvalidVariableName != err {
			t.Fatalf("Expected err: %v,\ngot: %v\n", localize.ErrInvalidVariableName, err)
		}
	}
}

// TestReservedVariableName ensures that reserved variable
// names cannot be assigned to the localized data.
func TestReservedVariableName(t *testing.T) {
	reservedCases := map[string]testCase{
		"var":      testCase{data: localize.Data{}},
		"function": testCase{data: localize.Data{}},
		"await":    testCase{data: localize.Data{}},
		"import":   testCase{data: localize.Data{}},
	}
	for name, tCase := range reservedCases {
		if _, err := localize.NewMap(name, tCase.data); localize.ErrReservedKeyword != err {
			t.Fatalf("Expected err: %v,\ngot: %v\n", localize.ErrReservedKeyword, err)
		}
	}
}
