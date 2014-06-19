package motto

import (
    "github.com/robertkrimen/otto"
)

type Motto interface {
    otto.Otto
    modules map[string]*Module
}

func (otto *Otto) RunModule(moduleFile string) (Value, error) {
    
}

type Module struct {
    Id string
    Filename string
    Loaded bool
    Parent *Module
    Children []*Module
}

func NewModule(moduleFile string) *Module {
    return &Module {
        Id: moduleFile,
        Filename: moduleFile,
    }
}

func Run(moduleFile string) (*Motto, otto.Value, error) {
    otto := New()
    value, err := otto.RunModule(moduleFile)

    return otto, value, err
}

