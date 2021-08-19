package main

import (
	"errors"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func init() {
	for k, v := range sprig.TxtFuncMap() {
		FunctionMap[k] = v
	}
}

var FunctionMap = template.FuncMap{
	"Nil": func() interface{} {
		return nil
	},
	"exit": func(message string) (string, error) {
		return "", errors.New(message)
	},
}
