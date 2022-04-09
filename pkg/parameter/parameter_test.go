package parameter_test

import (
	"testing"

	"github.com/balesz/protoc-gen-tmpl/pkg/parameter"
	"github.com/stretchr/testify/suite"
)

func TestParseSuite(t *testing.T) { suite.Run(t, new(ParseSuite)) }

type ParseSuite struct {
	suite.Suite
	parameter string
}

func (suite *ParseSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestWithEmptyParameter":
		suite.parameter = ""
	case "TestWithInvalidParameter":
		suite.parameter = ",,"
	case "TestWithOneParameter":
		suite.parameter = "param1=param1"
	case "TestWithTwoParameters":
		suite.parameter = "param1=param1,param2"
	}
}

func (suite *ParseSuite) TestWithEmptyParameter() {
	params, err := parameter.Parse(suite.parameter)
	suite.NoError(err)
	suite.Len(params, 0)
}

func (suite *ParseSuite) TestWithInvalidParameter() {
	params, err := parameter.Parse(suite.parameter)
	suite.NoError(err)
	suite.Len(params, 0)
}

func (suite *ParseSuite) TestWithOneParameter() {
	params, err := parameter.Parse(suite.parameter)
	suite.NoError(err)
	suite.Len(params, 1)
	suite.Equal(params["param1"], "param1")
}

func (suite *ParseSuite) TestWithTwoParameters() {
	params, err := parameter.Parse(suite.parameter)
	suite.NoError(err)
	suite.Len(params, 2)
	suite.Equal(params["param1"], "param1")
	suite.Equal(params["param2"], "true")
}

func TestBoolSuite(t *testing.T) { suite.Run(t, new(BoolSuite)) }

type BoolSuite struct{ suite.Suite }

func (suite *BoolSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestWithEmptyParameter":
		parameter.Parse("")
	case "TestWithoutValue":
		parameter.Parse("test")
	case "TestWithValue":
		parameter.Parse("test=false")
	}
}

func (suite *BoolSuite) TestWithEmptyParameter() {
	res := parameter.Bool("test", false)
	suite.False(res)
}

func (suite *BoolSuite) TestWithoutValue() {
	res := parameter.Bool("test", false)
	suite.True(res)
}

func (suite *BoolSuite) TestWithValue() {
	res := parameter.Bool("test", true)
	suite.False(res)
}

func TestIntSuite(t *testing.T) { suite.Run(t, new(IntSuite)) }

type IntSuite struct{ suite.Suite }

func (suite *IntSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestWithEmptyParameter":
		parameter.Parse("")
	case "TestWithoutValue":
		parameter.Parse("test=")
	case "TestWithValue":
		parameter.Parse("test=123")
	}
}

func (suite *IntSuite) TestWithEmptyParameter() {
	res := parameter.Int("test", 0)
	suite.Equal(res, 0)
}

func (suite *IntSuite) TestWithoutValue() {
	res := parameter.Int("test", 0)
	suite.Equal(res, 0)
}

func (suite *IntSuite) TestWithValue() {
	res := parameter.Int("test", 0)
	suite.Equal(res, 123)
}

func TestStringSuite(t *testing.T) { suite.Run(t, new(StringSuite)) }

type StringSuite struct{ suite.Suite }

func (suite *StringSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestWithEmptyParameter":
		parameter.Parse("")
	case "TestWithoutValue":
		parameter.Parse("test=")
	case "TestWithValue":
		parameter.Parse("test=hello world")
	}
}

func (suite *StringSuite) TestWithEmptyParameter() {
	res := parameter.String("test", "asdf")
	suite.Equal(res, "asdf")
}

func (suite *StringSuite) TestWithoutValue() {
	res := parameter.String("test", "asdf")
	suite.Empty(res)
}

func (suite *StringSuite) TestWithValue() {
	res := parameter.String("test", "asdf")
	suite.Equal(res, "hello world")
}
