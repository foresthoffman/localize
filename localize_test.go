/**
 * localize_test.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package localize

import (
	"html/template"
	"testing"
)

type testCase struct {
	data     Data
	expected []template.JS
}

var testCases = map[string]testCase{
	// Int case.
	"intCase": testCase{
		data: Data{
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
		data: Data{
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
		data: Data{
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
		data: Data{
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
var maps = make(map[string]*Map)

// TestNewMap insures that the data maps can be instantiated as
// expected.
func TestNewMap(t *testing.T) {
	for name, tCase := range testCases {
		m, err := NewMap(name, tCase.data)
		if nil != err {
			t.Fatalf("Failed to create new map (%q), err: %v\n", name, err)
		}
		maps[name] = m
	}
}

// TestJS insures that the localized data is formatted as
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
			t.Fatalf("Expected one of: %q, got: %q\n", testCases[name].expected, output)
		}
	}
}

// TestInvalidVariableName insures that invalid variable names
// cannot be assigned to the localized data.
func TestInvalidVariableName(t *testing.T) {
	invalidCases := map[string]testCase{
		"2var": testCase{data: Data{}},
		"-var": testCase{data: Data{}},
		"*var": testCase{data: Data{}},
	}
	for name, tCase := range invalidCases {
		if _, err := NewMap(name, tCase.data); ErrInvalidVariableName != err {
			t.Fatalf("Expected err: %v, got: %v\n", ErrInvalidVariableName, err)
		}
	}
}

// TestReservedVariableName insures that reserved variable
// names cannot be assigned to the localized data.
func TestReservedVariableName(t *testing.T) {
	reservedCases := map[string]testCase{
		"var":      testCase{data: Data{}},
		"function": testCase{data: Data{}},
		"await":    testCase{data: Data{}},
		"import":   testCase{data: Data{}},
	}
	for name, tCase := range reservedCases {
		if _, err := NewMap(name, tCase.data); ErrReservedKeyword != err {
			t.Fatalf("Expected err: %v, got: %v\n", ErrReservedKeyword, err)
		}
	}
}
