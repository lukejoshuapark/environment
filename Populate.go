package environment

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var configParsers = map[reflect.Type]reflect.Value{}
var stringType = reflect.TypeOf("")
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func UseParser(parserFunc interface{}) {
	t := reflect.TypeOf(parserFunc)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf(`cannot use "%v" as a parser function`, t))
	}

	if t.NumIn() != 1 || t.In(0) != stringType {
		panic("parser functions must accept a single string as input")
	}

	if t.NumOut() < 1 || t.NumOut() > 2 {
		panic("parser functions must return either T, or (T, error), where T can be any type but error")
	}

	if t.NumOut() == 2 && t.Out(1) != errorType {
		panic("parser functions must return either T, or (T, error), where T can be any type but error")
	}

	if t.Out(0) == errorType {
		panic("parser functions must return either T, or (T, error), where T can be any type but error")
	}

	configParsers[t.Out(0)] = reflect.ValueOf(parserFunc)
}

func Populate(into interface{}) error {
	tPtr := reflect.TypeOf(into)
	if tPtr.Kind() != reflect.Ptr {
		return fmt.Errorf(`expected Populate to be called on a pointer type, not "%v"`, tPtr)
	}

	v := reflect.ValueOf(into).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, ok := f.Tag.Lookup("environment")
		if !ok {
			continue
		}

		spl := strings.Split(tag, ",")
		if len(spl) < 1 || len(spl) > 2 {
			continue
		}

		varName := strings.TrimSpace(spl[0])
		varType := f.Type
		varValue, ok := os.LookupEnv(varName)
		if !ok {
			if len(spl) < 2 {
				return fmt.Errorf(`required environment variable "%v" for field "%v" on "%v" was missing`, varName, f.Name, t)
			}

			varValue = strings.TrimSpace(spl[1])
		}

		if varType == stringType {
			v.Field(i).Set(reflect.ValueOf(varValue))
			continue
		}

		parser, ok := configParsers[varType]
		if !ok {
			return fmt.Errorf(`the environment variable "%v" for field "%v" on "%v" was read, but no parser could be found for type "%v"`, varName, f.Name, t, varType)
		}

		result := parser.Call([]reflect.Value{reflect.ValueOf(varValue)})
		if len(result) == 2 && !result[1].IsNil() {
			err := result[1].Interface().(error)
			return fmt.Errorf(`the environment variable "%v" for field "%v" on "%v" was read, but the parser for type "%v" failed: %w`, varName, f.Name, t, varType, err)
		}

		v.Field(i).Set(result[0])
	}

	return nil
}
