package confik

import (
	"fmt"
	"strings"
)

type tag struct {
	Name     string
	Flags    []string
	Settings map[string]string
}

// ConfigTag represents the name, flags and settings on the struct field.
type ConfigTag struct {
	Name      string  // name of the environment variable
	Validator *string // field validator name
	Optional  bool    // is the environment variable optional?
	Default   *string // default value to use if the environment variable does not exist
	Unset     bool    // clear the environment variable after load?
}

// NewConfigTag will create a new [ConfigTag] with the default values.
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

func parseTag(tagStr string) (*tag, error) {
	if len(tagStr) == 0 {
		return nil, fmt.Errorf("invalid tag: empty tag")
	}
	parts := strings.Split(tagStr, ",")
	var tag tag
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
		// TODO remove in Go 1.22
		settingValue := settingValue
		switch settingName {
		case "validate":
			configTag.Validator = &settingValue
		case "default":
			configTag.Default = &settingValue
		default:
			return nil, fmt.Errorf("invalid env tag: unknown setting %s", settingName)
		}
	}

	return &configTag, nil
}
