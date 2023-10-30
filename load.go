// [confik] will build structs from environment files and variables.
//
// [confik] works by reading special struct tags on your fields.
//
// # Supported Types
//
//   - uint
//   - uint8
//   - uint16
//   - uint32
//   - uint64
//   - int
//   - int8
//   - int16
//   - int32
//   - int64
//   - float32
//   - float64
//   - bool
//   - [time.Duration]
//   - [time.Time]
//   - [url.URL]
//
// # Tag Options
//
// You can control how each field is configured through struct tags under the "env" key.
//
// Tags are specified like such:
//
//	type MyStruct struct {
//	  Name: `env:"NAME_OF_VARIABLE,flag1,flag2,setting1=value,setting2=value"`
//	}
//
// Available flags:
//
//   - optional: Dont require this value to exist in the environment.
//   - unset: Remove this environment value after load.
//
// Available settings:
//
//   - default=value: Set the default (string) value if it is not found in the environment.
//   - validator=validator: Set the name of the validator to use for this field.
//
// # Validators
//
// Fields can have their string values validated during environment load.
//
// Available validators:
//
//   - file: Verify that the path exists and is a file.
//   - dir: Verify that the path exists and is a directory.
//   - uri: Verify that the value is a URI.
//   - ip: Verify that the value is an IP address.
//   - port: Verify that the value is a port.
//   - hostport: Verify that the value is a host/port combination.
//   - cidr: Verify that the value is a CIDR.
//
// # Custom Validators
//
// Fields can be implement custom validators by specifying a [Validator] in [Config].
//
// See the examples below.
//
// # Custom Types
//
// Custom types can be supported by specifying a [Parser] in [Config].
//
// See the examples below.
//
// # Examples
package confik

import (
	"fmt"
	"os"
	"reflect"
)

// LoadFromEnv will build a T by reading values from environment files and variables.
func LoadFromEnv[T any](cfgs ...Config[T]) (*T, error) {
	cfg := DefaultConfig[T]()
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	// attempt to find and load the ".env" file
	if cfg.UseEnvFile {
		_, err := loadEnvFile(cfg)
		if err != nil {
			return nil, err
		}
	}

	var z T
	var t = reflect.TypeOf(z)

	// iterate over all the visible fields on the struct
	for _, field := range reflect.VisibleFields(t) {
		fieldConfig, err := newFieldConfig(cfg, field)
		if err != nil {
			return nil, err
		}

		// get a reflected value of the field
		var rv = reflect.ValueOf(&z).Elem().FieldByName(field.Name)

		// get the environment variable
		fieldValue, exists := os.LookupEnv(fieldConfig.Name)

		// unset the environment variable if applicable
		if fieldConfig.Unset {
			os.Unsetenv(fieldConfig.Name)
		}

		// handle default values if applicable
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

		// handle simple types and slices
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

		// handle more complex types (like time.Time, time.Duration, custom types)
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
