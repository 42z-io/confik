package confik

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// FieldConfig is the representation of the configuration for a field within a struct (after tags have been parsed).
type FieldConfig struct {
	ConfigTag            // the configuration specified in the tag
	Validate  *Validator // the custom validator for this field
}

// merge two maps together into a new map.
func mergeMap[T comparable, V any](a map[T]V, b map[T]V) map[T]V {
	merged := make(map[T]V)
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

// newFieldConfig will create a new FieldConfig for the given [reflect.StructField].
func newFieldConfig[T any](cfg Config[T], rv reflect.StructField) (*FieldConfig, error) {
	var fieldConfig FieldConfig
	tagStr := rv.Tag.Get("env")
	if tagStr != "" {
		tag, err := parseEnvTag(tagStr)
		if err != nil {
			return nil, fmt.Errorf("invalid tag on field %s: %w", rv.Name, err)
		}
		fieldConfig.ConfigTag = *tag
	} else {
		fieldConfig.ConfigTag = ConfigTag{
			Name: toEnvName(rv.Name),
		}
	}

	validators := mergeMap(fieldValidators, cfg.Validators)

	validatorName := fieldConfig.Validator
	if validatorName != nil {
		validator, exists := validators[*validatorName]
		if !exists {
			return nil, fmt.Errorf("unknown validator: %s", *validatorName)
		}
		fieldConfig.Validate = &validator
	}
	return &fieldConfig, nil
}

// verifyEnvName will ensure that a variable is in a suitable format for an environment variable.
func verifyEnvName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("invalid environment variable name: %s must be [A-Z0-9_]+", name)
	}
	for i := 0; i < len(name); i++ {
		if !(name[i] == '_' || (name[i] >= 'A' && name[i] <= 'Z') || (name[i] >= '0' && name[i] <= '9')) {
			return fmt.Errorf("invalid environment variable name: %s must be [A-Z0-9_]+", name)
		}
	}
	return nil
}

// toEnvName will take a field name and convert it into a format sutiable for an environment variable.
func toEnvName(name string) string {
	// split at capitalization, case change, or numbers
	var sb strings.Builder
	for i, c := range name {
		if i != 0 && (unicode.IsUpper(c) || unicode.IsDigit(c)) {
			sb.WriteString("_")
		}
		sb.WriteRune(unicode.ToUpper(c))
	}
	var res = sb.String()
	return res
}
