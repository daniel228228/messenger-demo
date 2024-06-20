package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadConfig(config any) error {
	v := reflect.ValueOf(config)

	if v.Kind() != reflect.Pointer || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("%q is not a struct pointer", v.Type().Name()))
	}

	walk(v.Elem())

	return nil
}

func walk(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)

		for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
			if value.IsNil() {
				panic(fmt.Sprintf("nil field %q", v.Type().Field(i).Name))
			} else {
				value = value.Elem()
			}
		}

		if value.Kind() == reflect.Struct {
			walk(value)
		} else {
			set(v.Type().Field(i), value)
		}
	}
}

func set(t reflect.StructField, v reflect.Value) {
	tag, ok := t.Tag.Lookup("env_config")
	if !ok {
		panic(fmt.Sprintf("config struct field %q has no \"env_config\" tag", t.Name))
	}

	value, ok := os.LookupEnv(tag)
	if !ok {
		panic(fmt.Sprintf("env variable %q does not exist", tag))
	}

	if !v.CanSet() {
		panic(fmt.Sprintf("can't access config struct field %q", t.Name))
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Bool:
		if val, err := strconv.ParseBool(value); err == nil {
			v.SetBool(val)
		} else {
			panic(err)
		}
	case reflect.Int:
		if val, err := strconv.ParseInt(value, 10, 0); err == nil {
			v.SetInt(val)
		} else {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("%s it not a supported config struct field type", v.Kind().String()))
	}
}
