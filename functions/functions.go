package functions

import (
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/balesz/protoc-gen-tmpl/data"
	"github.com/iancoleman/strcase"
)

func New(data *data.Data) *Functions {
	return &Functions{data: data, store: make(map[string]interface{})}
}

type Functions struct {
	data    *data.Data
	funcMap template.FuncMap
	store   map[string]interface{}
}

func (it *Functions) Map() template.FuncMap {
	if len(it.funcMap) != 0 {
		return it.funcMap
	}
	it.funcMap = template.FuncMap{
		"Nil":    it.nilFunc,
		"Exit":   it.exitFunc,
		"Fail":   it.failFunc,
		"ToList": it.toListFunc,

		"Set": it.setFunc,
		"Get": it.getFunc,

		"ToCamel":              strcase.ToCamel,
		"ToDelimited":          strcase.ToDelimited,
		"ToKebab":              strcase.ToKebab,
		"ToLowerCamel":         strcase.ToLowerCamel,
		"ToScreamingDelimited": strcase.ToScreamingDelimited,
		"ToScreamingKebab":     strcase.ToScreamingKebab,
		"ToScreamingSnake":     strcase.ToScreamingSnake,
		"ToSnake":              strcase.ToSnake,
		"ToSnakeWithIgnore":    strcase.ToSnakeWithIgnore,

		"FindFileByPath":       it.findFileByPathFunc,
		"FindDescriptorByName": it.findDescriptorByNameFunc,
		"FindServiceByName":    it.findServiceByNameFunc,
		"FindMessageByName":    it.findMessageByNameFunc,
		"FindEnumByName":       it.findEnumByNameFunc,

		"LeadingComments":         it.leadingCommentsFunc,
		"LeadingDetachedComments": it.leadingDetachedCommentsFunc,
		"TrailingComments":        it.trailingCommentsFunc,

		"Options": it.optionsFunc,
	}
	for k, v := range sprig.TxtFuncMap() {
		it.funcMap[k] = v
	}
	return it.funcMap
}

func (it *Functions) LookupExit(err error) (string, bool) {
	rxExit := regexp.MustCompile(".* error calling (Exit: .*)")
	if rxExit.MatchString(err.Error()) {
		matches := rxExit.FindStringSubmatch(err.Error())
		return matches[1], true
	}
	return "", false
}

func (it *Functions) LookupFail(err error) (string, bool) {
	rxFail := regexp.MustCompile(".* error calling (Fail: .*)")
	if rxFail.MatchString(err.Error()) {
		matches := rxFail.FindStringSubmatch(err.Error())
		return matches[1], true
	}
	return "", false
}

func (it *Functions) ResetStore() {
	it.store = make(map[string]interface{})
}
