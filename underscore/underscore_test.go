package underscore

import (
    "testing"
    "github.com/ddliu/motto"
)

func TestUnderscoreImport(t *testing.T) {
    _, v, _ := motto.Run("tests/index.js")
    i, _ := v.ToInteger()

    if i != 1 {
        t.Error("import underscore test failed")
    }
}