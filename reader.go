package confik

import (
	"fmt"
	"os"
	"reflect"
)

func LoadFromEnv[T any](cfgs ...Config) (*T, error) {
	cfg := DefaultConfig
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	// attempt to find and load the ".env" file
	if cfg.UseEnvFile {
		_, err := LoadEnvFile(cfg)
		if err != nil {
			return nil, err
		}
	}

	var z T
	var t = reflect.TypeOf(z)

	// iterate over all the visible fields on the struct
	for _, field := range reflect.VisibleFields(t) {
		fieldCfg, err := NewFieldConfig(cfg, field.Name, string(field.Tag))
		if err != nil {
			return nil, err
		}
		// get the environment variable
		fieldValue, exists := os.LookupEnv(fieldCfg.Name)

		// return an error if the environment variable doesn't exist and this field is not optional
		if !fieldCfg.Optional && !exists {
			return nil, fmt.Errorf("environment variable %s does not exist", fieldCfg.Name)
		}

		// skip to the next field if we cant find the environment variable
		if !exists {
			continue
		}

		// run validation on the environment variable (if any)
		if fieldCfg.Validator != nil {
			if err = fieldCfg.Validator(fieldCfg.Name, fieldValue); err != nil {
				return nil, err
			}
		}

		// get a reflected value of the field
		var rv = reflect.ValueOf(&z).Elem().FieldByName(field.Name)
		var kind = rv.Kind()
		if kind == reflect.Slice {
			err := handleSlice(fieldCfg, fieldValue, rv)
			if err != nil {
				return nil, err
			}
		} else {
			// check and see if there is a converter for kind of value this field has
			converter, converterExists := kindConverters[kind]
			customConverter, customConverterExists := cfg.CustomConverters[rv.Type()]
			if customConverter != nil {
				converter = customConverter
			}

			// if there is no standard converter, and no user provided convert return an error
			if !converterExists && !customConverterExists {
				return nil, fmt.Errorf("field %s of type %s is not supported", field.Name, field.Type)
			}

			// convert the value from a string to the fields type
			if err = converter(fieldCfg, fieldValue, rv); err != nil {
				return nil, err
			}
		}
	}
	return &z, nil
}
