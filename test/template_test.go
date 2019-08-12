/**
 * template_test.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestTemplate ensures that implementations of localize work
// as expected with the html.template package.
func TestTemplate(t *testing.T) {
	port := 3000
	dur := 2 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	go ListenAndServeWithClose(ctx, port)

	resp, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/")
	if nil != err {
		t.Fatalf("Failed to GET localhost address,\nerr: %v\n", err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if nil != err {
		t.Fatalf("Failed to read from response body,\nerr: %v\n", err)
	}
	// The order of elements in a map is not guaranteed,
	// therefore, both potential orders have to be checked.
	expected := []string{
		`_localData = {
"motd": [
"Hello world, welcome to a new day!",
],
"nonce": {
"login":"LaKJIIjIOUhjbKHdBJHGkhg",

},

};`,
		`_localData = {
"nonce": {
"login":"LaKJIIjIOUhjbKHdBJHGkhg",

},
"motd": [
"Hello world, welcome to a new day!",
],

};`,
	}
	matched := false
	for _, str := range expected {
		if index := strings.Index(string(contents), str); -1 != index {
			matched = true
		}
	}
	if !matched {
		t.Fatalf("Failed to find localized data,\nexpected one of: %v,\nbody: %v\n", expected, string(contents))
	}
	cancel()
}
