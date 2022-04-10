package config_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/balesz/protoc-gen-tmpl/internal/config"
	"github.com/balesz/protoc-gen-tmpl/internal/log"
	"github.com/stretchr/testify/suite"
)

func init() {
	log.Init(bytes.NewBufferString(""))
}

func TestLoad(t *testing.T) { suite.Run(t, new(LoadSuite)) }

type LoadSuite struct {
	suite.Suite
	configPath string
}

func (suite *LoadSuite) SetupSuite() {
	if dir, err := os.MkdirTemp("", "test-protoc-gen-tmpl-*"); err != nil {
		panic(err)
	} else {
		suite.configPath = dir
	}
}

func (suite *LoadSuite) TearDownSuite() {
	if err := os.RemoveAll(suite.configPath); err != nil {
		panic(err)
	}
}

func (suite *LoadSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestWithReadError":
		writeString(suite.configPath, "")
	case "TestWithParseError":
		writeString(suite.configPath, "hello: world")
	case "TestSuccessful":
		writeString(suite.configPath, "exclude: hello world")
	}
}

func (suite *LoadSuite) TestWithReadError() {
	res, err := config.Load("")
	suite.ErrorIs(err, config.ErrorReadFile)
	suite.Nil(res)
}

func (suite *LoadSuite) TestWithParseError() {
	res, err := config.Load(suite.configPath)
	suite.ErrorIs(err, config.ErrorInvalidFile)
	suite.Nil(res)
}

func (suite *LoadSuite) TestSuccessful() {
	res, err := config.Load(suite.configPath)
	suite.NoError(err)
	suite.NotNil(res)
	suite.Contains(res.Exclude, "hello world")
}

////////////////////////////////////////////////////////////////////////////////

func writeString(dir string, str string) {
	file, err := os.Create(path.Join(dir, "protoc-gen-tmpl.yaml"))
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if _, err := file.WriteString(str); err != nil {
		panic(err)
	} else if err := file.Sync(); err != nil {
		panic(err)
	}
}
