package confik

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type FieldConfig struct {
	ConfigTag
	Validate *Validator
}

func mergeMap(a map[string]Validator, b map[string]Validator) map[string]Validator {
	merged := make(map[string]Validator)
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

func NewFieldConfig[T any](cfg Config[T], rv reflect.StructField) (*FieldConfig, error) {
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

	validatorName := fieldConfig.ConfigTag.Validator
	if validatorName != nil {
		validator, exists := validators[*validatorName]
		if !exists {
			return &fieldConfig, nil
		}
		fieldConfig.Validate = &validator
	}
	return &fieldConfig, nil
}

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
