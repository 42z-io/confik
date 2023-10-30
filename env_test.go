package confik

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnvironmentVariableExpression(t *testing.T) {
	name, value, err := parseEnvironmentVariableExpression("VARIABLE=")
	assert.Nil(t, err)
	assert.Equal(t, "VARIABLE", name)
	assert.Equal(t, "", value)
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

func TestLoadEnvFileNoEnvFile(t *testing.T) {
	cwd, _ := os.Getwd()
	os.Chdir("/")
	defer func() {
		os.Chdir(cwd)
	}()

	kv, err := LoadEnvFile(Config[testAllTypes]{})
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{}, kv)
}

func TestLoadEnvFileDoesNotExist(t *testing.T) {
	_, err := LoadEnvFile(Config[testAllTypes]{
		EnvFilePath: ".fake",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment file does not exist: .fake", err.Error())
	}
}

func TestLoadEnvFileInvalid(t *testing.T) {
	_, err := LoadEnvFile(Config[testAllTypes]{
		EnvFilePath: "testdata/.invalid",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "invalid expression in env file: INVALID", err.Error())
	}
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
	assert.Nil(t, err)
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
		assert.Equal(t, "invalid expression in env file: BLAH", err.Error())
	}
}

func BenchmarkParseEnvFile(b *testing.B) {
	input := `
// Comment
TC_UNQUOTED=1
TC_QUOTED="hello world"
TC_SPACES=hello world
TC_LONG_NAME_1=test
# Comment


`
	scanner := bufio.NewScanner(strings.NewReader(input))
	for n := 0; n < b.N; n++ {
		_, err := ParseEnvFile(scanner)
		if err != nil {
			panic(err)
		}
	}
}
