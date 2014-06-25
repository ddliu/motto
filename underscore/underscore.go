package underscore

import (
    "github.com/ddliu/motto"
)

motto.Register("underscore", func(*motto.Motto) otto.Value)

motto.RegisterModule("underscore", underscore)

func underscore(vm *motto.Motto) (otto.Value, error) {
    
}