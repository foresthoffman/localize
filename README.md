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
    // Generates a new localization map with the provided data.
    dataMap, err := localize.NewMap(
        // This will tell the localizer to assign the data to
        // the "_localData" global JavaScript variable.
        "_localData",
        localize.Data{
            "motd": "Hello world, welcome to a new day!",

            // "nonce" will hold an object with an element with
            // the key, "login", and the value,
            // "LaKJIIjIOUhjbKHdBJHGkhg"
            "nonce": map[string]string{
                "login": "LaKJIIjIOUhjbKHdBJHGkhg",
            },
        },
    )

    // ...proper error handling, data manipulation, etc.
}
```

For a more complex example using the standard html/template and net/http packages check the [`test/template.go`](https://github.com/foresthoffman/localize/blob/master/test/template.go) file.

### How exactly are Golang data types translated to JavaScript?

`localize.ReflectTarget()` is the function that handles this process. From the docs:

> [ReflectTarget takes] a reflect.Value object and recursively determines the values of all the fields, sub-fields, elements, etc. At each step, the target's type is analyzed to see whether or not it's an enclosing type. If the target is an enclosing type, then the contents of the target will be wrapped appropriately. Square-brackets ("[]") are used for translating data to a JavaScript array. Curly-brackets ("{}") are used for translating data to a JavaScript object. Non-enclosing types [e.g. int, float64, and string] simply output according to their JavaScript equivalent. The complete contents of the top-most target is written piece-by-piece to the buffer provided.

Localized Golang interfaces are translated to JavaScript as native objects. As in, they are of type `Object`, but do not carry over their interface's context. So, if an instance of a struct `Page` were placed in a `localize.Map`, the JavaScript equivalent would not provide any explicit indication that the `Object` was a `Page`.

_That's all, enjoy!_
