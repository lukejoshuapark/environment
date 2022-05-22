package environment

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/lukejoshuapark/test"
	"github.com/lukejoshuapark/test/is"
)

func TestPopulateBasic(t *testing.T) {
	// Arrange.
	configParsers = map[reflect.Type]reflect.Value{}
	UseParser(ParseUInt16)
	os.Setenv("SOME_STRING_VALUE", "1234")
	os.Setenv("SOME_UINT16", "5678")

	cfg := &BasicConfig{}

	// Act.
	err := Populate(cfg)

	// Assert.
	test.That(t, err, is.Nil)
	test.That(t, cfg.SomeStringValue, is.EqualTo("1234"))
	test.That(t, cfg.SomeUInt16, is.EqualTo(uint16(5678)))
	test.That(t, cfg.SomeOptionalValue, is.EqualTo("9012"))
}

func TestPopulateRequiredVariableMissing(t *testing.T) {
	// Arrange.
	configParsers = map[reflect.Type]reflect.Value{}
	UseParser(ParseUInt16)
	os.Unsetenv("SOME_STRING_VALUE")
	os.Setenv("SOME_UINT16", "5678")

	cfg := &BasicConfig{}

	// Act.
	err := Populate(cfg)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, err.Error(), is.EqualTo(`required environment variable "SOME_STRING_VALUE" for field "SomeStringValue" on "environment.BasicConfig" was missing`))
}

func TestPopulateMissingParser(t *testing.T) {
	// Arrange.
	configParsers = map[reflect.Type]reflect.Value{}
	os.Setenv("SOME_STRING_VALUE", "1234")
	os.Setenv("SOME_UINT16", "5678")

	cfg := &BasicConfig{}

	// Act.
	err := Populate(cfg)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, err.Error(), is.EqualTo(`the environment variable "SOME_UINT16" for field "SomeUInt16" on "environment.BasicConfig" was read, but no parser could be found for type "uint16"`))
}

func TestPopulateParserFailure(t *testing.T) {
	// Arrange.
	configParsers = map[reflect.Type]reflect.Value{}
	UseParser(ParseUInt16)
	os.Setenv("SOME_STRING_VALUE", "1234")
	os.Setenv("SOME_UINT16", "567891")

	cfg := &BasicConfig{}

	// Act.
	err := Populate(cfg)

	// Assert.
	test.That(t, err, is.NotNil)
	test.That(t, err.Error(), is.EqualTo(`the environment variable "SOME_UINT16" for field "SomeUInt16" on "environment.BasicConfig" was read, but the parser for type "uint16" failed: strconv.ParseUint: parsing "567891": value out of range`))
}

// --

func ParseUInt16(val string) (uint16, error) {
	v, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(v), nil
}

// --

type BasicConfig struct {
	SomeStringValue   string `environment:"SOME_STRING_VALUE"`
	SomeUInt16        uint16 `environment:"SOME_UINT16"`
	SomeOptionalValue string `environment:"SOME_OPTIONAL_VALUE,9012"`
}
