## Localize

Package localize provides functions for translating Golang data structures to JavaScript primitives. The translated, or "localized", JavaScript that can be produced by this package is intended to be used directly with the html/template package. This package eases the process of passing global data down to front-end scripts.

### Installation

Run `go get github.com/foresthoffman/localize`

### Importing

Import the package by including `github.com/foresthoffman/localize` in your import block.

e.g.

```Go
package main

import(
	...
	"github.com/foresthoffman/localize"
)
```

### Usage

Here's a simple example of the syntax:

```Go
import(
    "github.com/foresthoffman/localize"
)

func main() {

    // generates a new localization map with the provided data
    dataMap, err := localize.NewMap(
        &map[string]interface{}{

            // "motd" will hold an array with an element with the value, "Hello world, welcome
            // to a new day!"
            "motd": "Hello world, welcome to a new day!",

            // "nonce" will hold an object with an element with the key, "login", and the value,
            // "LaKJIIjIOUhjbKHdBJHGkhg"
            "nonce": map[string]string{
                "login": "LaKJIIjIOUhjbKHdBJHGkhg",
            },
        },

    // this will tell the localizer to assign the data to the "_localData" global JavaScript
    // variable
    ).SetVarName("_localData")

    // ...proper error handling, data manipulation, etc.
}
```

Here's a more complex example using the standard html/template and net/http packages (this one is fully functional; copy-paste it, and try it out):

```Go
package main

import (
    "github.com/foresthoffman/localize"
    "html/template"
    "net/http"
)

var tmpl *template.Template
var page *Page

type Page struct {
    LocalizedData *localize.LocalizeMap
}

func RootHandler(w http.ResponseWriter, rq *http.Request) {

    // Executes the template, and runs any template actions, which includes the LocalizedData's
    // "JS()" function. The template will be returned to the client's browser along with the
    // new JavaScript data.
    err := tmpl.Execute(w, *page)
    if nil != err {
        panic(err)
    }
}

func main() {

    // prepares the localized data
    dataMap, err := localize.NewMap(
        &map[string]interface{}{

            // "motd" will hold an array with an element with the value, "Hello world, welcome
            // to a new day!"
            "motd": "Hello world, welcome to a new day!",

            // "nonce" will hold an object with an element with the key, "login", and the value,
            // "LaKJIIjIOUhjbKHdBJHGkhg"
            "nonce": map[string]string{
                "login": "LaKJIIjIOUhjbKHdBJHGkhg",
            },
        },
    ).SetVarName("_localData")
    if nil != err {
        panic(err)
    }

    // sets up a page that will provide the template with the LocalizedData field
    page = &Page{
        LocalizedData: dataMap,
    }

    // normally this would be in an HTML file on its own, but for the sake of brevity...
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

                    // Access the first element of the motd property of the _localData variable
                    // to get the message of the day, and then insert it into the motd span of
                    // the header tag on the page.
                    document.querySelector(".page .motd").innerText = _localData.motd[0];
                };
            </script>
        </body>
        </html>
    `
    tmpl, err = template.New("hello").Parse(templateBody)
    if nil != err {
        panic(err)
    }

    // This fires up the webserver and waits for connections to "http://localhost:3000/". Hitting
    // that page will present the client with the following header text:
    //
    // "The message of the day is: Hello world, welcome to a new day!"
    http.HandleFunc("/", RootHandler)
    http.ListenAndServe(":3000", nil)
}
```

### How exactly are Golang data types translated to JavaScript?

`localize.ReflectTarget()` is the function that handles this process. From the docs:

> [ReflectTarget takes] a reflect.Value object and recursively determines the values of all the fields, sub-fields, elements, etc. At each step, the target's type is analyzed to see whether or not it's an enclosing type. If the target is an enclosing type, then the contents of the target will be wrapped appropriately. Square-brackets ("[]") are used for translating data to a JavaScript array. Curly-brackets ("{}") are used for translating data to a JavaScript object. Non-enclosing types [e.g. int, float64, and string] simply output according to their JavaScript equivalent. The complete contents of the top-most target is written piece-by-piece to the buffer provided.

Localized Golang interfaces are translated to JavaScript as native objects. As in, they are of type `Object`, but do not carry over their interface's context. So, if an instance of a struct `Page` were placed in a `localize.LocalizeMap`, the JavaScript equivalent would not provide any explicit indication that the `Object` was a `Page`.

_That's all, enjoy!_
