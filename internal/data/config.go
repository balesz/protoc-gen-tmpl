package data

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func InitConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("protoc-gen-tmpl")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

type Config struct {
	Exclude []string
	Output  []struct {
		Name string
		Path string
	}
	Types struct {
		Enum     string
		Map      string
		Message  string
		Repeated string
		Scalar   struct {
			Double   string
			Float    string
			Int32    string
			Int64    string
			Uint32   string
			Uint64   string
			Sint32   string
			Sint64   string
			Fixed32  string
			Fixed64  string
			Sfixed32 string
			Sfixed64 string
			Bool     string
			String   string
			Bytes    string
		}
	}
}

func (cfg *Config) IsExcluded(name string) bool {
	for _, pattern := range cfg.Exclude {
		if regexp.MustCompile(pattern).MatchString(name) {
			return true
		}
	}
	return false
}

func (cfg *Config) OutputByName(name string) string {
	for _, out := range cfg.Output {
		if strings.HasSuffix(name, out.Name) {
			return out.Path
		}
	}
	return ""
}

func (cfg *Config) TypeOf(field protoreflect.FieldDescriptor) string {
	var result string

	if field.IsMap() {
		result = strings.ReplaceAll(cfg.Types.Map, "{{TKey}}", cfg.TypeOf(field.MapKey()))
		result = strings.ReplaceAll(result, "{{TValue}}", cfg.TypeOf(field.MapValue()))
		return result
	}

	switch field.Kind() {
	case protoreflect.EnumKind:
		result = parse(cfg.Types.Enum, field.Enum())
	case protoreflect.MessageKind:
		result = parse(cfg.Types.Message, field.Message())
	case protoreflect.DoubleKind:
		result = cfg.Types.Scalar.Double
	case protoreflect.FloatKind:
		result = cfg.Types.Scalar.Float
	case protoreflect.Int32Kind:
		result = cfg.Types.Scalar.Int32
	case protoreflect.Int64Kind:
		result = cfg.Types.Scalar.Int64
	case protoreflect.Uint32Kind:
		result = cfg.Types.Scalar.Uint32
	case protoreflect.Uint64Kind:
		result = cfg.Types.Scalar.Uint64
	case protoreflect.Sint32Kind:
		result = cfg.Types.Scalar.Sint32
	case protoreflect.Fixed32Kind:
		result = cfg.Types.Scalar.Fixed32
	case protoreflect.Fixed64Kind:
		result = cfg.Types.Scalar.Fixed64
	case protoreflect.Sfixed32Kind:
		result = cfg.Types.Scalar.Sfixed32
	case protoreflect.BoolKind:
		result = cfg.Types.Scalar.Bool
	case protoreflect.StringKind:
		result = cfg.Types.Scalar.String
	case protoreflect.BytesKind:
		result = cfg.Types.Scalar.Bytes
	}

	if field.IsList() {
		result = strings.Replace(cfg.Types.Repeated, "{{T}}", result, 1)
	}

	return result
}

func parse(format string, descriptor protoreflect.Descriptor) string {
	funcMap := template.FuncMap{
		"ToCamel":              strcase.ToCamel,
		"ToDelimited":          strcase.ToDelimited,
		"ToKebab":              strcase.ToKebab,
		"ToLowerCamel":         strcase.ToLowerCamel,
		"ToScreamingDelimited": strcase.ToScreamingDelimited,
		"ToScreamingKebab":     strcase.ToScreamingKebab,
		"ToScreamingSnake":     strcase.ToScreamingSnake,
		"ToSnake":              strcase.ToSnake,
		"ToSnakeWithIgnore":    strcase.ToSnakeWithIgnore,
	}
	buf := new(bytes.Buffer)
	if tmpl, err := template.New("type").Funcs(funcMap).Parse(format); err != nil {
		panic(err)
	} else if err := tmpl.Execute(buf, descriptor); err != nil {
		panic(err)
	} else {
		return buf.String()
	}
}
