# motto

[![Build Status](https://travis-ci.org/ddliu/motto.png)](https://travis-ci.org/ddliu/motto)

Modular [otto](https://github.com/robertkrimen/otto)

Motto provide a module environment to run javascript file in golang.

## Installation

```bash
go get github.com/ddliu/motto
```

## Usage

A javascript module:

```js
// content of module1.js

function test() {
    console.log("test");
}

exports.test = test;

// Module export works like Nodejs

// export with module.exports
// module.exports.test = test

// export with "return"
// return test
```

The main js file:

```js
// content of index.js

var _ = require('path/to/underscore');
var module1 = require('path/to/module1.js')

module1.test();

console.log(_.min([3,2,1,4,5]);
```

Run index.js:

```go
package main

import (
    "github.com/ddliu/motto"
)

func main() {
    motto.Run("path/to/index.js")
}
```

You can also install the motto command line tool to run it directly:

```bash
go install github.com/ddliu/motto/motto
motto path/to/index.js
```

## TODO

- More tests

## Changelog

### v0.1.0 (2014-06-22)

Initial release