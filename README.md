# motto

[![Build Status](https://travis-ci.org/ddliu/motto.png)](https://travis-ci.org/ddliu/motto)

Modular [otto](https://github.com/robertkrimen/otto)

Motto provide a Nodejs like module environment to run javascript files in golang.

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

## Nodejs Module Capable

The module environment is almost capable with Nodejs [spec](http://nodejs.org/api/modules.html).

Some Nodejs modules(without core module dependencies) can be used directly in Motto.

## Create Core Module

You can implement a Nodejs module in Motto and use it in your javascript file.

Refer to the test file for more details.

## TODO

- More tests

## Changelog

### v0.1.0 (2014-06-22)

Initial release

### v0.2.0 (2014-06-24)

Make module capable with Nodejs