package confik

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
func TestFindEnvFile(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() {
		os.Chdir(cwd)
	}()
	os.Chdir("testdata/folder1/folder2/folder3")
	path, err := FindEnvFile()
	assert.Nil(t, err)
	absPath := filepath.Join(cwd, "testdata/folder1/.env")
	assert.Equal(t, absPath, path)
}

func TestFindEnvFileNotFound(t *testing.T) {
	cwd, _ := os.Getwd()
	defer func() {
		os.Chdir(cwd)
	}()
	os.Chdir("/")
	path, err := FindEnvFile()
	assert.Nil(t, err)
	assert.Equal(t, "", path)
}

func TestParseEnvFile(t *testing.T) {
	input := `
// Comment
TC_UNQUOTED=1
TC_QUOTED="hello world"
TC_SPACES=hello world
TC_LONG_NAME_1=test
# Comment


`
	scanner := bufio.NewScanner(strings.NewReader(input))
	kv, err := ParseEnvFile(scanner)
	assert.Nil(t, err, "expected parseEnvFile to not fail")
	expected := map[string]string{
		"TC_UNQUOTED":    "1",
		"TC_QUOTED":      "hello world",
		"TC_SPACES":      "hello world",
		"TC_LONG_NAME_1": "test",
	}
	assert.Equal(t, len(expected), len(kv))
	for testK, testV := range expected {
		value, exists := kv[testK]
		assert.True(t, exists, "expected to find key %s", testK)
		assert.Equal(t, testV, value, "invalid value for key %s", testK)
	}
}

func TestParseEnvFileInvalid(t *testing.T) {
	input := "BLAH"
	scanner := bufio.NewScanner(strings.NewReader(input))
	_, err := ParseEnvFile(scanner)
	if assert.Error(t, err) {
		assert.Equal(t, "invalid line in env file: BLAH", err.Error())
	}
}
