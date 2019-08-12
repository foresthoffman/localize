/**
 * template.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

package test

import (
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/foresthoffman/localize"
)

var tmpl *template.Template
var page *Page

// Page is a wrapper for the localized data, to be used with a
// HTML template.
type Page struct {
	LocalizedData *localize.Map
}

// RootHandler executes the template, and runs any template
// actions, which includes the LocalizedData's "JS()" function.
// The template will be returned to the client's browser along
// with the new JavaScript data.
func RootHandler(w http.ResponseWriter, rq *http.Request) {
	err := tmpl.Execute(w, *page)
	if nil != err {
		panic(err)
	}
}

// ListenAndServeWithClose listens for connections while the
// context is alive.
func ListenAndServeWithClose(ctx context.Context, port int) error {
	// Prepares the localized data.
	dataMap, err := localize.NewMap(
		"_localData",
		localize.Data{
			"motd": "Hello world, welcome to a new day!",
			"nonce": map[string]string{
				"login": "LaKJIIjIOUhjbKHdBJHGkhg",
			},
		},
	)
	if nil != err {
		return err
	}

	// Sets up a page that will provide the template with the
	// LocalizedData field.
	page = &Page{
		LocalizedData: dataMap,
	}

	// Normally this would be in an HTML file on its own, but
	// for the sake of brevity...
	templateBody := `
        <!DOCTYPE html>
        <html>
        <head>
            <title>Hello world!</title>
        </head>
        <body>
            <div class="page">
                <h1>The message of the day is: <span class="motd"></span></h1>
            </div>

            <!--
            calls the "JS()" function of the "LocalizedData" of the
            object that was passed to the template.
            -->
            <script type="text/javascript">{{.LocalizedData.JS}}</script>
            <script type="text/javascript">
                window.onload = function() {

                    // Access the first element of the motd
                    // property of the _localData variable to
                    // get the message of the day, and then
                    // insert it into the motd span of the
                    // header tag on the page.
                    document.querySelector(".page .motd").innerText = _localData.motd[0];
                };
            </script>
        </body>
        </html>
    `
	tmpl, err = template.New("hello").Parse(templateBody)
	if nil != err {
		return err
	}

	// This fires up the webserver and waits for connections to
	// "http://localhost:3000/". Hitting that page will present
	// the client with the following header text:
	//
	// "The message of the day is: Hello world, welcome to a
	// new day!"
	http.HandleFunc("/", RootHandler)
	server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil}

	go server.ListenAndServe()
	<-ctx.Done()
	server.Close()

	return nil
}
