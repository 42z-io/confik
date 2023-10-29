package confik

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldNameToEnvironmentVariable(t *testing.T) {
	res := map[string]string{
		"test":           "TEST",
		"test3":          "TEST_3",
		"testCamelCase":  "TEST_CAMEL_CASE",
		"testCamel3Case": "TEST_CAMEL_3_CASE",
		"other_name":     "OTHER_NAME",
	}
	for input, expect := range res {
		result, err := fieldNameToEnvironmentVariable(input)
		assert.Nil(t, err, "expected field name %s to be valid", input)
		assert.Equal(t, expect, result)
	}
}
func TestValidateEnvironmentVariable(t *testing.T) {
	valid := []string{
		"TEST_1",
		"TEST_1_OTHER",
		"OTHEROTHEROTHER",
		"T",
		"TEST_1_2_3333333",
		"3_3_3",
		"33____3",
	}

	invalid := []string{
		"@",
		"",
		"MY-NAME",
		"PEICE3$",
		"^_3030",
	}

	for _, v := range valid {
		err := validateEnvironmentVariable(v)
		assert.Nil(t, err, "expected variable %s to be valid", v)
	}

	for _, i := range invalid {
		err := validateEnvironmentVariable(i)
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Sprintf("invalid environment variable name: %s must be [A-Z0-9_]+", i), err.Error(), "expected variable %s to be invalid", i)
		}
	}
}
