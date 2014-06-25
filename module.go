package motto


type ModuleInterface interface {
    GetValue() (otto.Value, error)
    GetId() string
}

type ModuleLoader struct {
    Id string,
    Loader 
}

// Node capable module implement, see: http://nodejs.org/api/modules.html
type Module struct {
    Id string
    Filename string
    Loaded bool
    vm *Motto
    // Parent *Module
    // Children []*Module
    Value otto.Value // module return value
}

// Load another module by name
func (this *Module) Require(name string) (otto.Value, error) {
    module, err := this.FindModule(name)

    if err != nil {
        return otto.UndefinedValue(), err
    }

    if this.vm.HasModule(module.GetId()) {
        module, _ = this.vm.GetModule(module.GetId())
        v, err := module.GetValue()
        if err != nil {
            return otto.UndefinedValue(), err
        }

        return v, nil
    }

    // new module
    this.vm.AddModule(module)

    return module.GetValue()
}

func (this *Module) GetValue() (otto.Value, error) {
    if this.Loaded {
        return this.Value, nil
    }

    // execute module
    
    moduleSource, err := ioutil.ReadFile(this.Filename)

    if err != nil {
        return otto.UndefinedValue(), err
    }

    // load json
    if filepath.Ext(this.Filename) == ".json" {
        value, err := this.vm.Call("JSON.parse", nil, string(moduleSource))
        if err != nil {
            return otto.UndefinedValue(), err
        }

        this.Value = value
        this.Loaded = true

        return this.Value, nil
    }

    // execute js
    moduleSource = append([]byte("(function(module) {var require = module.require;var exports = module.exports;\n"), moduleSource...)
    moduleSource = append(moduleSource, []byte("\n})")...)

    // Provide the "require" method in the module scope.
    jsRequire := func(call otto.FunctionCall) otto.Value {
        jsModuleFilename := call.Argument(0).String()

        moduleValue, err := this.Require(jsModuleFilename)
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
    this.Loaded = true
    this.Value = moduleValue

    return this.Value, nil
}

func (this *Module) GetId() string {
    return this.Id
}

// Find module by name
func (this *Module) FindModule(name string) (ModuleInterface, error) {
    var err error
    if len(name) == 0 {
        return nil, errors.New("Empty module name")
    }

    // paths to locate `name`
    var paths []string
    // path
    if name[0] == '.' || name[0] == '/' {
        if name[0] == '.' {
            name = filepath.Join(filepath.Dir(this.Filename), name)
        }
        if !filepath.IsAbs(name) {
            if name, err = filepath.Abs(name); err != nil {
                return nil, err
            }
        }
        paths = append(paths, name)
        ext := filepath.Ext(name)
        if ext != ".js" && ext != ".json" {
            paths = append(paths, name + ".js", name + ".json")
        }
    } else if module, ok := this.vm.GetModule(name); ok {
        return module, nil
    } else {
        // current_module/node_modules/xxx
        paths = append(paths, filepath.Join(filepath.Dir(this.Filename), "node_modules", name))

        // module paths registered in vm
        for _, v := range this.vm.paths {
            paths = append(paths, filepath.Join(v, name))
        }
    }


    for _, v := range paths {
        ok, err := isDir(v)
        if err != nil {
            return nil, err
        }

        if ok {
            packageJsonFilename := filepath.Join(v, "package.json")
            ok, err := isFile(packageJsonFilename)
            if err != nil {
                return nil, err
            }

            var entryPoint string
            if ok {
                entryPoint, err = parsePackageEntryPoint(packageJsonFilename)
                if err != nil {
                    return nil, err
                }
            } else {
                entryPoint = "./index.js"
            }

            return &Module {
                Id: filepath.Join(v, entryPoint),
                Filename: filepath.Join(v, entryPoint),
                vm: this.vm,
            }, nil
        }

        ok, err = isFile(v)
        if err != nil {
            return nil, err
        }

        if ok {
            return &Module {
                Id: v,
                Filename: v,
                vm: this.vm,
            }, nil
        }
    }

    return nil, errors.New("Module not found: " + name)
}