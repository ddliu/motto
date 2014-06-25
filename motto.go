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
    // "fmt"
)

type Motto struct {
    *otto.Otto
    modules map[string]ModuleInterface
    paths []string
}

// Run a js file as a module
func (this *Motto) RunModule(name string) (otto.Value, error) {

    // if name is a file, convert it to the absolute path.
    // Because it might not be recognized by Module.FindModule
    if ok, _ := isFile(name); ok {
        if absPath, err := filepath.Abs(name); err == nil {
            name = absPath
        }
    }
    baseModule := &Module {
        Id: ".",
        Filename: ".",
        vm: this,
    }

    return baseModule.Require(name)
}

// Get a registered module by id.
func (this *Motto) GetModule(id string) (ModuleInterface, bool) {
    module, ok := this.modules[id]
    return module, ok
}

// Check if specified module id exists
func (this *Motto) HasModule(id string) bool {
    _, ok := this.modules[id]

    return ok
}

func (this *Motto) FindModule(name string) (ModuleInterface, error) {
    baseModule := &Module {
        Id: "",
        Filename: "",
        vm: this,
    }

    return baseModule.FindModule(name)
}

// Add new modules to current vm.
func (this *Motto) AddModule(modules ...ModuleInterface) {
    if this.modules == nil {
        this.modules = make(map[string]ModuleInterface)
    }

    for _, module := range modules {
        this.modules[module.GetId()] = module
    }
}

func (this *Motto) AddPath(paths ...string) {
    this.paths = append(this.paths, paths...)
}

// Run module by name in the motto module environment.
func Run(name string) (*Motto, otto.Value, error) {
    vm := &Motto {otto.New(), nil, nil}
    v, err := vm.RunModule(name)

    return vm, v, err
}

func New() *Motto {
    return &Motto {otto.New(), nil, nil}
}