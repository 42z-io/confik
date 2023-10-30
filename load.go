package confik

import (
	"fmt"
	"os"
	"reflect"
)

func LoadFromEnv[T any](cfgs ...Config[T]) (*T, error) {
	cfg := DefaultConfig[T]()
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
		fieldConfig, err := NewFieldConfig(cfg, field)
		if err != nil {
			return nil, err
		}

		// get a reflected value of the field
		var rv = reflect.ValueOf(&z).Elem().FieldByName(field.Name)

		// get the environment variable
		fieldValue, exists := os.LookupEnv(fieldConfig.Name)
		if !exists && cfg.DefaultValue != nil {
			var drv = reflect.ValueOf(cfg.DefaultValue).Elem().FieldByName(field.Name)
			rv.Set(drv)
			continue
		} else if !exists && fieldConfig.Default != nil {
			fieldValue = *fieldConfig.Default
			exists = true
		}

		// return an error if the environment variable doesn't exist and this field is not optional
		if !fieldConfig.Optional && !exists {
			return nil, fmt.Errorf("environment variable %s does not exist and has no default", fieldConfig.Name)
		}

		// skip to the next field if we cant find the environment variable
		if !exists {
			continue
		}

		// run validation on the environment variable (if any)
		if fieldConfig.Validate != nil {
			if err = (*fieldConfig.Validate)(fieldConfig.Name, fieldValue); err != nil {
				return nil, err
			}
		}

		var kind = rv.Kind()
		if kind == reflect.Slice {
			err := handleSlice(fieldConfig, fieldValue, rv)
			if err != nil {
				return nil, err
			}
			continue
		}

		kindParser, exists := kindParsers[kind]
		if exists {
			// convert the value from a string to the fields type
			if err = kindParser(fieldConfig, fieldValue, rv); err != nil {
				return nil, err
			}
			continue
		}

		parsers := mergeMap(typeParsers, cfg.Parsers)
		parser, exists := parsers[field.Type]
		if exists {
			// convert the value from a string to the fields type
			if err = parser(fieldConfig, fieldValue, rv); err != nil {
				return nil, err
			}
			continue
		}

		return nil, fmt.Errorf("field %s of type %s has no parser", field.Name, field.Type)
	}
	return &z, nil
}
