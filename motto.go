package motto

import (
    "github.com/robertkrimen/otto"
    "strings"
)

type Motto interface {
    otto.Otto
    modules map[string]*Module
    currentModule *Module
}

func (motto *Motto) Require(moduleFile string) (otto.Value, error) {

    // before
    // exec
    // after
}

func (motto *Motto) NewModule(moduleFile) (*Module) {
    return &Module {
        Id: moduleFile,
        Filename: moduleFile,
        vm: motto,
    }
}


type Module struct {
    Id string
    Filename string
    Loaded bool
    vm *Motto
    // Parent *Module
    // Children []*Module
    Value otto.Value // module return value
    ModuleValue otto.Object // the module variable in js
}

func (m *Module) ToJavascriptModule() {
    // module
    module := otto.Object {}
    
    // require
    require := func(call otto.FunctionCall) otto.Value {
        filename, err := call.Argument(0).ToString()
        if err != nil || filename == "" || !strings.HasSuffix(strings.ToLower(filename), ".js") {
            return otto.UndefinedValue()
        }

        currentModule = motto.currentModule

        if !filepath.IsAbs(filename) {
            filename, err = filepath.Abs(filepath.Join(filepath.Dir(currentModule.Filename), filename))
            if err != nil {
                return otto.UndefinedValue()
            }

            existModule, ok := motto.modules[filename]
            // exist
            if ok {
                // Circle reference
                if !existModule.Loaded {
                    panic(call.Otto.Call("new Error", nil, "Circle module reference detected"))
                } else {
                    return existModule.Value
                }
            }

        }

        f, err := ioutil.ReadFile(filename)
        if err == nil {
            panic(call.Otto.Call("new Error", nil, fmt.Sprintf("Module file cannot be read: %s", filename)))
        }

        v, err := call.Otto.Run(f)

        module := &Module {
            Id: filename,
            Filename: filename,
            Loaded: false,
        }

        motto.modules[filename] = module

        motto.currentModule = module


        call.Otto.run()


        currentModule.Filename
    }

    // exports
}

func (this *Module) Require(filename string) (otto.Value, error) {
    absModulePath, err := this.resolvePath(filename)
    if err != nil {
        return err
    }

    existModule, ok = this.vm.modules[absModulePath]
    if ok {
        if existModule.Loaded {
            return existModule.Value, nil
        } else {
            return nil, errors.New("Circle module dependencies detected")
        }
    }

    module = &Module {
        Id: absModulePath,
        Filename: absModulePath,
        vm: this.motto,
    }

    this.vm.modules[absModulePath] = module

    // execute module
    moduleSource, err := ioutil.ReadFile(absModulePath)

    if err != nil {
        return nil, err
    }

    moduleSource =  
    
    module.Loaded = true
    module.Value = value

    return module.Value, nil
}

func (this *Module) resolvePath(filename string) (string, error) {
    if !filepath.IsAbs(filename) {
        filename, err := filepath.Join(this.Filename, filename)
        if err != nil {
            return "", err
        }
    }

    return filename
}

func MainModule(file string) *Module {
    fileAbs, err := filepath.Abs(file)
    if err != nil {
        panic(fmt.Printf("Path %s can not be resolved"))
    }

    return NewModule(fileAbs)
}

func NewModule(moduleFile string) *Module {
    return &Module {
        Id: moduleFile,
        Filename: moduleFile,
    }
}

func Run(file string) (*Motto, otto.Value, error) {
    otto := New()
    value, err := otto.RunModule(moduleFile)

    return otto, value, err
}