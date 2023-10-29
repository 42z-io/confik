package confik

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fatih/structtag"
)

type FieldConfig struct {
	Optional  bool
	Name      string
	Validator FieldValidator
}

func NewFieldConfig(cfg Config, name string, tag string) (*FieldConfig, error) {
	envName, err := fieldNameToEnvironmentVariable(name)
	if err != nil {
		return nil, err
	}
	fieldCfg := FieldConfig{
		Name:      envName,
		Optional:  false,
		Validator: nil,
	}
	tags, err := structtag.Parse(tag)
	if err != nil {
		return nil, err
	}
	envTags, err := tags.Get("env")
	if err != nil {
		return &fieldCfg, nil
	}

	if err := validateEnvironmentVariable(envTags.Name); err != nil {
		return nil, fmt.Errorf("invalid struct tag on %s: %w", name, err)
	}

	fieldCfg.Name = envTags.Name

	if envTags.HasOption("optional") {
		fieldCfg.Optional = true
	}
	for k, v := range fieldValidators {
		if envTags.HasOption(k) {
			fieldCfg.Validator = v
		}
	}
	if cfg.CustomValidators != nil {
		for k, v := range cfg.CustomValidators {
			if envTags.HasOption(k) {
				fieldCfg.Validator = v
			}
		}
	}
	return &fieldCfg, nil
}

func validateEnvironmentVariable(name string) error {
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

func fieldNameToEnvironmentVariable(name string) (string, error) {
	// split at capitalization, case change, or numbers
	var sb strings.Builder
	for i, c := range name {
		if i != 0 && (unicode.IsUpper(c) || unicode.IsDigit(c)) {
			if _, err := sb.WriteString("_"); err != nil {
				return "", err
			}
		}
		sb.WriteRune(unicode.ToUpper(c))
	}
	var res = sb.String()
	return res, nil
}
