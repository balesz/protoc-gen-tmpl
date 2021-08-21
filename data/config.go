package data

import (
	"regexp"
	"strings"

	"github.com/spf13/viper"
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
