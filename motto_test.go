package motto

import (
    "testing"
    "github.com/robertkrimen/otto"
    "io/ioutil"
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

func TestNpmModule(t *testing.T) {
    _, v, err := Run("tests/npm/index.js")

    if err != nil {
        t.Error(err)
    }

    i, _ := v.ToInteger()

    if i != 1 {
        t.Error("npm test failed: ", i , "!=", 1)
    }
}

type testModuleFS struct {
    vm *Motto
}

func (this *testModuleFS) GetId() string {
    return "fs"
}

func (this *testModuleFS) GetValue() (otto.Value, error) {
    fs, _ := this.vm.Object(`({})`)
    fs.Set("readFileSync", this.readFileSync)
    return this.vm.ToValue(fs)
}

func (this *testModuleFS) readFileSync(call otto.FunctionCall) otto.Value {
    filename, _ := call.Argument(0).ToString()
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return otto.UndefinedValue()
    }

    v, _ := call.Otto.ToValue(string(bytes))
    return v
}


func TestCoreModule(t *testing.T) {
    vm := New()
    vm.AddModule(&testModuleFS{vm})

    v, err := vm.RunModule("tests/core_module_test.js")
    if err != nil {
        t.Error(err)
    }

    s, _ := v.ToString()
    if s != "cat" {
        t.Error("core module test failed: ", s, "!=", "cat")
    }
}