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

func (motto *Motto) RunModule(moduleFile string) (otto.Value, error) {
        
}

func (motto *Motto) Inject() {
    motto.Set("require", func(call otto.FunctionCall) otto.Value {
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




    })
}

type Module struct {
    Id string
    Filename string
    Loaded bool
    // Parent *Module
    // Children []*Module
    Value otto.Value
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