// Copyright 2014 dong<ddliuhb@gmail.com>.
// Licensed under the MIT license.
// 
// Motto - Modular otto
// 
// Motto wraps otto to provide a Nodejs like module environment.
package motto

import (
    "github.com/robertkrimen/otto"
    "errors"
    "io/ioutil"
    "path/filepath"
)

type Motto struct {
    *otto.Otto
    modules map[string]*Module
}

// Run a js file as a module
func (this *Motto) RunModule(filename string) (otto.Value, error) {
    module := &Module {
        Id: ".",
        Filename: ".",
        vm: this,
    }

    return module.Require(filename)
}

// Get a loaded module by filename.
func (this *Motto) GetModule(filename string) (*Module, bool) {
    module, ok := this.modules[filename]
    return module, ok
}

// Add a new module to current vm.
func (this *Motto) AddModule(filename string, module *Module) {
    if this.modules == nil {
        this.modules = make(map[string]*Module)
    }

    this.modules[filename] = module
}

func (this *Motto) ResolveModuleId(name string) (id string, filename string, err error) {
    m, ok := this.modules[name]
    if ok {
        return m.Id, m.Filename, nil
    }

    if name == "" {
        return "", "", errors.New("Empty module name")
    }

    // file path
    if name[0] == '.' || name[0] == '/' {
        filename, err := filepath.Abs(name)
        if err != nil {
            return "", "", err
        }

        ok, err := isDir(filename)
        if err != nil {
            return "", "", err
        }

        if ok {
            // has package.json
            packageJsonFilename := filepath.Join(filename, "package.json")
            ok, err := isFile(packageJsonFilename)

            if err != nil {
                return "", "", err
            }

            if ok {
                index, err := parsePackageJsonIndex(packageJsonFilename)
                if err != nil {
                    return "", "", err
                }

                return filename, filepath.Join(filename, index), nil
            }

            // todo
            // return filename, 

        }

        return filename, filename, nil
    }

    return "", "", nil
}

type Module struct {
    Id string
    Filename string
    Loaded bool
    vm *Motto
    // Parent *Module
    // Children []*Module
    Value otto.Value // module return value
}

// Load another module by filename
func (this *Module) Require(filename string) (otto.Value, error) {
    absModulePath, err := this.resolvePath(filename)
    if err != nil {
        return otto.UndefinedValue(), err
    }

    existModule, ok := this.vm.GetModule(absModulePath)
    if ok {
        if existModule.Loaded {
            return existModule.Value, nil
        } else {
            return otto.UndefinedValue(), errors.New("Circle module dependencies detected")
        }
    }

    module := &Module {
        Id: absModulePath,
        Filename: absModulePath,
        vm: this.vm,
    }

    this.vm.AddModule(absModulePath, module)

    // execute module
    moduleSource, err := ioutil.ReadFile(absModulePath)

    if err != nil {
        return otto.UndefinedValue(), err
    }

    moduleSource = append([]byte("(function(module) {var require = module.require;var exports = module.exports;\n"), moduleSource...)
    moduleSource = append(moduleSource, []byte("\n})")...)

    // Provide the "require" method in the module scope.
    jsRequire := func(call otto.FunctionCall) otto.Value {
        jsModuleFilename := call.Argument(0).String()

        moduleValue, err := module.Require(jsModuleFilename)
        if err != nil {
            jsException(this.vm, "Error", "motto: " + err.Error())
        }

        return moduleValue
    }

    jsModule, _ := this.vm.Object(`({exports: {}})`)
    jsModule.Set("require", jsRequire)
    jsExports, _ := jsModule.Get("exports")

    // Run the module source, with "jsModule" as the "module" varaible, "jsExports" as "this"(Nodejs capable).
    moduleReturn, err := this.vm.Call(string(moduleSource), jsExports, jsModule)
    if err != nil {
        return otto.UndefinedValue(), err
    }

    var moduleValue otto.Value
    if !moduleReturn.IsUndefined() {
        moduleValue = moduleReturn
        jsModule.Set("exports", moduleValue)
    } else {
        moduleValue, _ = jsModule.Get("exports")
    }
    module.Loaded = true
    module.Value = moduleValue

    return module.Value, nil
}

// Get absolute path of another module
func (this *Module) resolvePath(filename string) (string, error) {
    var err error
    if !filepath.IsAbs(filename) {
        filename = filepath.Join(filepath.Dir(this.Filename), filename)
        if !filepath.IsAbs(filename) {
            filename, err = filepath.Abs(filename)
            if err != nil {
                return "", err
            }
        }
    }

    return filename, nil
}

// Run javascript file in the motto module environment.
func Run(filename string) (*Motto, otto.Value, error) {
    vm := &Motto {otto.New(), nil}
    v, err := vm.RunModule(filename)

    return vm, v, err
}

// Throw a javascript error, see https://github.com/robertkrimen/otto/issues/17
func jsException(vm *Motto, errorType, msg string) {
    value, _ := vm.Call("new " + errorType, nil, msg)
    panic(value)
}
