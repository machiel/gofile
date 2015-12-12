# Overview
[![Twitter](https://img.shields.io/badge/author-%40MachielMolenaar-blue.svg)](https://twitter.com/MachielMolenaar)
[![GoDoc](https://godoc.org/github.com/Machiel/gofile?status.svg)](https://godoc.org/github.com/Machiel/gofile)

Gofile is a library that abstracts away the filesystem operations for you, and
allows you to plug-in your own filesystem implementation.

Currently there are two drivers:

* godropbox - Dropbox implementation
* golocal - Local FS implementation

Please mind, this is a very alpha release, so please double check the code
before trying it in production. Heck, it doesn't even have tests yet :).
Be warned.

Gofile was greatly inspired by (copied from)
[@frankdejonge](https://twitter.com/frankdejonge)'s Flysystem.

# License
Gofile is licensed under a MIT license.

# Installation
To get started with using gofile, you require gofile, and at least one driver.

In order to install gofile: `go get github.com/Machiel/gofile`

If you want to get started using the local driver you'll also have to install
that one:
`go get github.com/Machiel/gofile/golocal`

If you want to use the Dropbox driver, you'll have to install some dependencies:

```
go get golang.org/x/oauth2
go get github.com/stacktic/dropbox
```

And then get godropbox:

`go get github.com/Machiel/gofile/godropbox`

# Usage

## Example
```go
package main

import (
    "github.com/Machiel/gofile"
    _ "github.com/Machiel/gofile/golocal"
    _ "github.com/Machiel/gofile/godropbox"
)


func main() {
    // godropbox initialization
    fs, err := gofile.New("dropbox", map[string]string{
        "client_id":     "",
        "client_secret": "",
        "token":         "",
    })

    handleError(err)

    err = fs.Write("myfile.txt", "Hello world!")

    handleError(err)

    // golocal initialization
    localFs, err := gofile.New("local", map[string]string{
        "rootDir" : "/tmp",
    })

    handleError(err)
}
```

# Contributions
Contributions are more than welcome, you can contribute by providing more
implementations, like AWS S3 or FTP.

Other than that, as mentioned earlier, currently there are 0 tests written,
so if you want you can help out with that stuff :)!

Other examples:
* Return streams instead of data arrays
* Visibility settings (public/private)
