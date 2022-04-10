package config

import (
	"errors"
	"fmt"

	"github.com/balesz/protoc-gen-tmpl/internal/log"
	"github.com/spf13/viper"
)

var ErrorReadFile = errors.New("can't read file")

var ErrorInvalidFile = errors.New("the file has wrong format")

func Load(path string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("protoc-gen-tmpl")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		log.Error("viper.ReadInConfig: %v", err)
		return nil, fmt.Errorf("config.Load: %w", ErrorReadFile)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		log.Error("viper.Unmarshal: %v", err)
		return nil, fmt.Errorf("config.Load: %w", ErrorInvalidFile)
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
