package motto

import (
    "testing"
)

func TestModule(t *testing.T) {
    _, v, err := Run("tests/index.js")
    if err != nil {
        t.Error(err)
    }

    i, _ := v.ToString()
    if i != "rat" {
        t.Error("testing result: ", i , "!=", "rat")
    }
}