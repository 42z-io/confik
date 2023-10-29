package confik

import (
	"fmt"
	"os"
	"reflect"
)

func NewEnvReader[T any](cfgs ...Config) (*T, error) {
	cfg := DefaultConfig
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	// Attempt to find and load the .env file
	if cfg.UseEnvFile {
		_, err := LoadEnvFile(cfg)
		if err != nil {
			return nil, err
		}
	}

	var z T
	var t = reflect.TypeOf(z)
	for _, field := range reflect.VisibleFields(t) {
		fieldCfg, err := NewFieldConfig(cfg, field.Name, string(field.Tag))
		if err != nil {
			return nil, err
		}
		fieldValue, exists := os.LookupEnv(fieldCfg.Name)
		if !fieldCfg.Optional && !exists {
			return nil, fmt.Errorf("environment variable %s does not exist", fieldCfg.Name)
		}
		if !exists {
			continue
		}
		if fieldCfg.Validator != nil {
			if err = fieldCfg.Validator(fieldCfg.Name, fieldValue); err != nil {
				return nil, err
			}
		}
		var rv = reflect.ValueOf(&z).Elem().FieldByName(field.Name)
		var kind = rv.Kind()
		if kind == reflect.Slice {
			err := handleSlice(fieldCfg, fieldValue, rv)
			if err != nil {
				return nil, err
			}
		} else {
			converter, exists1 := typeConverters[kind]
			customConverter, exists2 := cfg.CustomConverters[rv.Type()]
			if customConverter != nil {
				converter = customConverter
			}
			if !exists1 && !exists2 {
				return nil, fmt.Errorf("field %s of type %s is not supported", field.Name, field.Type)
			}
			if err = converter(fieldCfg, fieldValue, rv); err != nil {
				return nil, err
			}
		}
	}
	return &z, nil
}
