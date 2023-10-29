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

type SimpleEnv struct {
	Hello string
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

type MyTestConfig struct {
	AStringList []string
	CustomName  string  `env:"A_CUSTOM_NAME"`
	AUint8      uint8   `env:"UINT8"`
	AUint16     uint16  `env:"UINT16"`
	AUint32     uint32  `env:"UINT32"`
	AUint64     uint64  `env:"UINT64"`
	Aint8       int8    `env:"INT8"`
	Aint16      int16   `env:"INT16"`
	Aint32      int32   `env:"INT32"`
	Aint64      int64   `env:"INT64"`
	AFloat32    float32 `env:"FLOAT32"`
	AFloat64    float64 `env:"FLOAT64"`
	ABool       bool    `env:"BOOL"`
	Optional    bool    `env:"OPTIONAL,optional"`
}

func TestNewEnvReader(t *testing.T) {
	os.Clearenv()
	os.Setenv("INT16", "42")
	cfg, err := NewEnvReader[MyTestConfig]()
	assert.Nil(t, err)
	assert.Equal(t, true, cfg.ABool)
	assert.Equal(t, "A custom name", cfg.CustomName)
	assert.Equal(t, []string{"ABC", "DEF", "GHI"}, cfg.AStringList)
	assert.Equal(t, uint8(255), cfg.AUint8)
	assert.Equal(t, uint16(65535), cfg.AUint16)
	assert.Equal(t, uint32(4294967295), cfg.AUint32)
	assert.Equal(t, uint64(18446744073709551615), cfg.AUint64)
	assert.Equal(t, int8(-128), cfg.Aint8)
	assert.Equal(t, int16(42), cfg.Aint16)
	assert.Equal(t, int32(-2147483648), cfg.Aint32)
	assert.Equal(t, int64(-9223372036854775808), cfg.Aint64)
	assert.Equal(t, float32(65.3918495928), cfg.AFloat32)
	assert.Equal(t, float64(3.14159), cfg.AFloat64)
}

type MyTestConfigWithValidatorAndOptional struct {
	Website  string `env:"WEBSITE,uri"`
	Optional bool   `env:"OPTIONAL,optional"`
}

func TestNewEnvReaderWithValidator(t *testing.T) {
	os.Clearenv()
	cfg, err := NewEnvReader[MyTestConfigWithValidatorAndOptional](Config{
		EnvFilePath: "testdata/.uri",
		UseEnvFile:  true,
	})
	assert.Nil(t, err)
	assert.Equal(t, "https://google.com/my_site", cfg.Website)
	assert.Equal(t, false, cfg.Optional)
}

type MyTestConfigInvalid struct {
	Invalid string `env:"@@,uri"`
}

func TestNewEnvReaderInvalidTags(t *testing.T) {
	_, err := NewEnvReader[MyTestConfigInvalid](Config{
		EnvFilePath: "testdata/.uri",
		UseEnvFile:  true,
	})
	if assert.Error(t, err) {
		assert.Equal(t, "invalid struct tag on Invalid: invalid environment variable name: @@ must be [A-Z0-9_]+", err.Error())
	}
}
func TestNewEnvReaderOverride(t *testing.T) {
	os.Clearenv()
	os.Setenv("INT16", "42")
	cfg, err := NewEnvReader[MyTestConfig](Config{
		UseEnvFile:      true,
		EnvFileOverride: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, int16(-32768), cfg.Aint16)
}

func TestNewEnvReaderMissing(t *testing.T) {
	os.Clearenv()
	_, err := NewEnvReader[MyTestConfig](Config{
		UseEnvFile:      false,
		EnvFileOverride: false,
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment variable A_STRING_LIST does not exist", err.Error())
	}
}

func TestNewEnvReaderEnvNotExist(t *testing.T) {
	os.Clearenv()
	_, err := NewEnvReader[MyTestConfig](Config{
		UseEnvFile:  true,
		EnvFilePath: ".fake",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment file does not exist: .fake", err.Error())
	}
}

func TestNewEnvReaderEnvIsDir(t *testing.T) {
	os.Clearenv()
	_, err := NewEnvReader[MyTestConfig](Config{
		UseEnvFile:  true,
		EnvFilePath: "testdata/",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment file is a directory: testdata/", err.Error())
	}
}

func TestNewEnvReaderFoundEnvIsDir(t *testing.T) {
	os.Clearenv()
	cwd, _ := os.Getwd()
	target := "testdata/folder1/folder2/folder3/folder4"
	os.Chdir(target)
	defer func() {
		os.Chdir(cwd)
	}()
	_, err := NewEnvReader[MyTestConfig]()
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Sprintf("environment file is a directory: %s", filepath.Join(cwd, target, ".env")), err.Error())
	}
}
