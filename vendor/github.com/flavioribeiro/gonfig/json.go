package gonfig

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Gonfig implementation
// Implements the Gonfig interface
type JsonGonfig struct {
	obj map[string]interface{}
}

// FromJsonFile opens the file supplied and calls
// FromJson function
func FromJsonFile(filename string) (Gonfig, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config, err := FromJson(f)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// FromJson reads the contents from the supplied reader.
// The content is parsed as json into a map[string]interface{}.
// It returns a JsonGonfig struct pointer and any error encountered
func FromJson(reader io.Reader) (Gonfig, error) {
	jsonBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		return nil, err
	}
	return &JsonGonfig{obj}, nil
}

// GetString uses Get to fetch the value behind the supplied key.
// It returns a string with either the retreived value or the default value and any error encountered.
// If value is not a string it returns a UnexpectedValueTypeError
func (jgonfig *JsonGonfig) GetString(key string, defaultValue interface{}) (string, error) {
	configValue, err := jgonfig.Get(key, defaultValue)
	if err != nil {
		return "", err
	}
	if stringValue, ok := configValue.(string); ok {
		return stringValue, nil
	} else {
		return "", &UnexpectedValueTypeError{key: key, value: configValue, message: "value is not a string"}
	}
}

// GetInt uses Get to fetch the value behind the supplied key.
// It returns a int with either the retreived value or the default value and any error encountered.
// If value is not a int it returns a UnexpectedValueTypeError
func (jgonfig *JsonGonfig) GetInt(key string, defaultValue interface{}) (int, error) {
	value, err := jgonfig.GetFloat(key, defaultValue)
	if err != nil {
		return -1, err
	}
	return int(value), nil
}

// GetFloat uses Get to fetch the value behind the supplied key.
// It returns a float with either the retreived value or the default value and any error encountered.
// It returns a bool with either the retreived value or the default value and any error encountered.
// If value is not a float it returns a UnexpectedValueTypeError
func (jgonfig *JsonGonfig) GetFloat(key string, defaultValue interface{}) (float64, error) {
	configValue, err := jgonfig.Get(key, defaultValue)
	if err != nil {
		return -1.0, err
	}
	if floatValue, ok := configValue.(float64); ok {
		return floatValue, nil
	} else if intValue, ok := configValue.(int); ok {
		return float64(intValue), nil
	} else {
		return -1.0, &UnexpectedValueTypeError{key: key, value: configValue, message: "value is not a float"}
	}
}

// GetBool uses Get to fetch the value behind the supplied key.
// It returns a bool with either the retreived value or the default value and any error encountered.
// If value is not a bool it returns a UnexpectedValueTypeError
func (jgonfig *JsonGonfig) GetBool(key string, defaultValue interface{}) (bool, error) {
	configValue, err := jgonfig.Get(key, defaultValue)
	if err != nil {
		return false, err
	}
	if boolValue, ok := configValue.(bool); ok {
		return boolValue, nil
	} else {
		return false, &UnexpectedValueTypeError{key: key, value: configValue, message: "value is not a bool"}
	}
}

// GetAs uses Get to fetch the value behind the supplied key.
// The value is serialized into json and deserialized into the supplied target interface.
// It returns any error encountered.
func (jgonfig *JsonGonfig) GetAs(key string, target interface{}) error {
	configValue, err := jgonfig.Get(key, nil)
	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(configValue)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return err
	}
	return nil
}

// Get attempts to retreive the value behind the supplied key.
// It returns a interface{} with either the retreived value or the default value and any error encountered.
// If supplied key is not found and defaultValue is set to nil it returns a KeyNotFoundError
// If supplied key path goes deeper into a non-map type (string, int, bool) it returns a UnexpectedValueTypeError
func (jgonfig *JsonGonfig) Get(key string, defaultValue interface{}) (interface{}, error) {
	parts := strings.Split(key, "/")
	var tmp interface{} = jgonfig.obj
	for index, part := range parts {
		if len(part) == 0 {
			continue
		}
		if confMap, ok := tmp.(map[string]interface{}); ok {
			if value, exists := confMap[part]; exists {
				tmp = value
			} else if defaultValue != nil {
				return defaultValue, nil
			} else {
				return nil, &KeyNotFoundError{key: path.Join(append(parts[:index], part)...)}
			}
		} else {
			return nil, &UnexpectedValueTypeError{key: path.Join(parts[:index]...), value: tmp, message: "value behind key is not a map[string]interface{}"}
		}
	}
	return tmp, nil
}
