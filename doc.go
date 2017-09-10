/**
 * doc.go
 *
 * Copyright (c) 2017 Forest Hoffman. All Rights Reserved.
 * License: MIT License (see the included LICENSE file) or download at
 *     https://raw.githubusercontent.com/foresthoffman/localize/master/LICENSE
 */

/*
Package localize provides functions for translating Golang data structures to JavaScript primitives.
The translated, or "localized", JavaScript that can be produced by this package is intended to be
used directly with the "html/template" package. This package eases the process of passing global data
down to front-end scripts.

Here's an example:

    import(
        "html/template"
        "github.com/foresthoffman/localize"
    )

    func main() {
        ...coming soon
    }
*/

package localize
