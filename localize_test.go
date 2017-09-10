/**
 * localize_test.go
 *
 * Copyright (c) 2017 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package localize

import (
	"html/template"
	"testing"
)

// Tests that the JavaScript block returned by the localize.JS() function is as expected.
func TestJS(t *testing.T) {

	// int case
	intMap, err := NewMap(
		&map[string]interface{}{
			"int": 1954,
		},
	).SetVarName("_localIntData")
	if nil != err {
		t.Fatal(err)
	}

	intJS := intMap.JS()
	expectedintJS := template.JS(
		`_localIntData = {
"int": [
1954,
],

};`,
	)
	if expectedintJS != intJS {
		t.Errorf(
			"int output mismatch:\ngot\n%v\n, want\n%v\n",
			intJS,
			expectedintJS,
		)
	}
	// int case end

	// intArray case
	intArrayMap := NewMap(
		&map[string]interface{}{
			"intArray": []int{1, 2, 3, 4, 5},
		},
	)
	intArrayJS := intArrayMap.JS()
	expectedIntArrayJS := template.JS(
		`_globalVars = {
"intArray": [
[1,2,3,4,5,],

],

};`,
	)
	if expectedIntArrayJS != intArrayJS {
		t.Errorf(
			"intArray output mismatch:\ngot\n%v\n, want\n%v\n",
			intArrayJS,
			expectedIntArrayJS,
		)
	}
	// intArray case end

	// arrayArray case
	arrayArrayMap := NewMap(
		&map[string]interface{}{
			"arrayArray": [][]int{
				[]int{6, 7, 8, 9, 10},
				[]int{11, 12, 13, 14, 15},
			},
		},
	)
	arrayArrayJS := arrayArrayMap.JS()
	expectedArrayArrayJS := template.JS(
		`_globalVars = {
"arrayArray": [
[[6,7,8,9,10,],
[11,12,13,14,15,],
],

],

};`,
	)
	if expectedArrayArrayJS != arrayArrayJS {
		t.Errorf(
			"arrayArray output mismatch:\ngot\n%v\n, want\n%v\n",
			arrayArrayJS,
			expectedArrayArrayJS,
		)
	}
	// arrayArray case end

	// assocArray case
	assocArrayMap := NewMap(
		&map[string]interface{}{
			"assocArray": map[string]string{
				"baz": "fubar",
				"foo": "bar",
			},
		},
	)
	assocArrayJS := assocArrayMap.JS()
	expectedassocArrayJSA := template.JS(
		`_globalVars = {
"assocArray": {
"baz":"fubar",
"foo":"bar",

},

};`,
	)
	expectedassocArrayJSB := template.JS(
		`_globalVars = {
"assocArray": {
"foo":"bar",
"baz":"fubar",

},

};`,
	)
	if expectedassocArrayJSA != assocArrayJS && expectedassocArrayJSB != assocArrayJS {
		t.Errorf(
			"assocArray output mismatch:\ngot\n%v\n, want\n%v, or\n%v\n",
			assocArrayJS,
			expectedassocArrayJSA,
			expectedassocArrayJSB,
		)
	}
	// assocArray case end
}
