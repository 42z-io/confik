package confik

import (
	"fmt"
	"strings"
)

type Tag struct {
	Name     string
	Flags    []string
	Settings map[string]string
}

type ConfigTag struct {
	Name      string
	Validator *string
	Optional  bool
	Default   *string
	Unset     bool
}

func NewConfigTag(name string) ConfigTag {
	return ConfigTag{
		Name:      name,
		Validator: nil,
		Optional:  false,
		Default:   nil,
		Unset:     false,
	}
}

func parseSetting(expressionStr string) (string, string, error) {
	parts := strings.Split(expressionStr, "=")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid setting %s: invalid syntax", expressionStr)
	}
	return parts[0], parts[1], nil
}

func parseTag(tagStr string) (*Tag, error) {
	if len(tagStr) == 0 {
		return nil, fmt.Errorf("empty tag")
	}
	parts := strings.Split(tagStr, ",")
	var tag Tag
	tag.Name = parts[0]
	tag.Settings = make(map[string]string)
	tag.Flags = make([]string, 0)
	for _, expressionStr := range parts[1:] {
		if strings.Contains(expressionStr, "=") {
			settingName, settingValue, err := parseSetting(expressionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid tag: %w", err)
			}
			tag.Settings[settingName] = settingValue
		} else {
			tag.Flags = append(tag.Flags, expressionStr)
		}
	}
	return &tag, nil

}

func parseEnvTag(tagStr string) (*ConfigTag, error) {
	tag, err := parseTag(tagStr)
	if err != nil {
		return nil, fmt.Errorf("invalid env tag: %w", err)
	}

	if err := verifyEnvName(tag.Name); err != nil {
		return nil, fmt.Errorf("invalid env tag: %w", err)
	}

	configTag := NewConfigTag(tag.Name)
	for _, flagName := range tag.Flags {
		switch flagName {
		case "optional":
			configTag.Optional = true
		case "unset":
			configTag.Unset = true
		default:
			return nil, fmt.Errorf("invalid env tag: unknown flag %s", flagName)
		}
	}

	for settingName, settingValue := range tag.Settings {
		var settingValueC = strings.Clone(settingValue)
		switch settingName {
		case "validate":
			configTag.Validator = &settingValueC
		case "default":
			configTag.Default = &settingValueC
		default:
			return nil, fmt.Errorf("invalid env tag: unknown setting %s", settingName)
		}
	}

	return &configTag, nil
}
