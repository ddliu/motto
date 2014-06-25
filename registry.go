package motto

var globalModules map[string]AddonModule = make(map[string]ModuleLoader)

type ModuleLoader func(*Motto) (otto.Value, error)

func RegisterModule(id string, m ModuleLoader) {
    globalModules[id] = m
}