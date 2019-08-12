/**
 * doc.go
 *
 * Copyright (c) 2017-2019 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

/*
Package localize provides functions for translating Golang data
structures to JavaScript primitives. The translated, or
"localized", JavaScript that can be produced by this package is
intended to be used directly with the html/template package.
This package eases the process of passing global data down to
front-end scripts.

Here's a simple example of the syntax:

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

For a more complex example using the standard html/template and
net/http packages check the test/template.go file.

*/
package localize
