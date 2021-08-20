package main

import (
	"errors"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	"findFileByPath": func(files []protoreflect.FileDescriptor, path string) protoreflect.FileDescriptor {
		for _, file := range files {
			if strings.HasSuffix(file.Path(), path) {
				return file
			}
		}
		return nil
	},
}
