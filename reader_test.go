package confik

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAllTypes struct {
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

func TestLoadFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("INT16", "42")
	cfg, err := LoadFromEnv[testAllTypes]()
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

type testCustomValidator struct {
	Website  string `env:"WEBSITE,uri"`
	Optional bool   `env:"OPTIONAL,optional"`
}

func TestLoadFromEnvWithValidator(t *testing.T) {
	os.Clearenv()
	cfg, err := LoadFromEnv[testCustomValidator](Config{
		EnvFilePath: "testdata/.uri",
		UseEnvFile:  true,
	})
	assert.Nil(t, err)
	assert.Equal(t, "https://google.com/my_site", cfg.Website)
	assert.Equal(t, false, cfg.Optional)
}

type testInvalidTags struct {
	Invalid string `env:"@@,uri"`
}

func TestLoadFromEnvInvalidTags(t *testing.T) {
	_, err := LoadFromEnv[testInvalidTags](Config{
		EnvFilePath: "testdata/.uri",
		UseEnvFile:  true,
	})
	if assert.Error(t, err) {
		assert.Equal(t, "invalid struct tag on Invalid: invalid environment variable name: @@ must be [A-Z0-9_]+", err.Error())
	}
}

type testUnsupportedSlice struct {
	Custom []MyCustomType
}

func TestLoadFromEnvUnsupportedSlice(t *testing.T) {
	os.Clearenv()
	os.Setenv("CUSTOM", "Hello")
	_, err := LoadFromEnv[testUnsupportedSlice](Config{
		UseEnvFile: false,
	})
	if assert.Error(t, err) {
		assert.Equal(t, "CUSTOM is invalid: []confik.MyCustomType is not supported", err.Error())
	}
}
func TestLoadFromEnvEnvFileOverride(t *testing.T) {
	os.Clearenv()
	os.Setenv("INT16", "42")
	cfg, err := LoadFromEnv[testAllTypes](Config{
		UseEnvFile:      true,
		EnvFileOverride: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, int16(-32768), cfg.Aint16)
}

func TestLoadFromEnvRequiredFields(t *testing.T) {
	os.Clearenv()
	_, err := LoadFromEnv[testAllTypes](Config{
		UseEnvFile:      false,
		EnvFileOverride: false,
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment variable A_STRING_LIST does not exist", err.Error())
	}
}

func TestLoadFromEnvEnvFileNotExist(t *testing.T) {
	os.Clearenv()
	_, err := LoadFromEnv[testAllTypes](Config{
		UseEnvFile:  true,
		EnvFilePath: ".fake",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment file does not exist: .fake", err.Error())
	}
}

func TestLoadFromEnvEnvFileIsDir(t *testing.T) {
	os.Clearenv()
	_, err := LoadFromEnv[testAllTypes](Config{
		UseEnvFile:  true,
		EnvFilePath: "testdata/",
	})
	if assert.Error(t, err) {
		assert.Equal(t, "environment file is a directory: testdata/", err.Error())
	}
}

func TestLoadFromEnvFoundEnvIsDir(t *testing.T) {
	os.Clearenv()
	cwd, _ := os.Getwd()
	target := "testdata/folder1/folder2/folder3/folder4"
	os.Chdir(target)
	defer func() {
		os.Chdir(cwd)
	}()
	_, err := LoadFromEnv[testAllTypes]()
	if assert.Error(t, err) {
		assert.Equal(t, fmt.Sprintf("environment file is a directory: %s", filepath.Join(cwd, target, ".env")), err.Error())
	}
}
