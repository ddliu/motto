package motto

import (
    "testing"
)

func TestModule(t *testing.T) {
    motto, v, err := Run("tests/index.js")
    if err != nil {
        t.Error(err)
    }

    // test v
}