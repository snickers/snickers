/*
Package gonfig implements methods for reading config files encoded in json using path expressions and default fallback values
*/
package gonfig

import (
	"fmt"
	"reflect"
)

type Gonfig interface {
	Get(key string, defaultValue interface{}) (interface{}, error)
	GetString(key string, defaultValue interface{}) (string, error)
	GetInt(key string, defaultValue interface{}) (int, error)
	GetFloat(key string, defaultValue interface{}) (float64, error)
	GetBool(key string, defaultValue interface{}) (bool, error)
	GetAs(key string, target interface{}) error
}

type KeyNotFoundError struct {
	key string
}

func (err *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key not found, key: %s", err.key)
}

type UnexpectedValueTypeError struct {
	key     string
	value   interface{}
	message string
}

func (err *UnexpectedValueTypeError) Error() string {
	return fmt.Sprintf("%s, key: %s, value: %v (%s)", err.message, err.key, err.value, reflect.TypeOf(err.value).Name())
}
