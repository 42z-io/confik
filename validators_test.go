package confik

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func customValidator(envName string, value string) error {
	return fmt.Errorf("%s-%s", envName, value)
}

type CustomValidatorType struct {
	MyField string `env:"MY_FIELD,validate=custom"`
}

func TestLoadFromEnvCustomValidator(t *testing.T) {
	os.Clearenv()
	os.Setenv("MY_FIELD", "test")
	_, err := LoadFromEnv(Config[CustomValidatorType]{
		UseEnvFile: false,
		Validators: map[string]Validator{
			"custom": customValidator,
		},
	})
	if assert.Error(t, err) {
		assert.Equal(t, "MY_FIELD-test", err.Error())
	}
}

func TestValidatorUri(t *testing.T) {
	err := validateUri("MY_VAR", "my_value")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=my_value is not a URI: parse \"my_value\": invalid URI for request", err.Error())
	}
	err = validateUri("MY_VAR", "https://google.com/my_path?hello=world")
	assert.Nil(t, err)
}

func TestValidatorIp(t *testing.T) {
	err := validateIp("MY_VAR", "128")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=128 is not a valid IP: invalid format", err.Error())
	}
	valids := []string{
		"192.168.0.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8:3c4d:15::1a2f:1a2b",
	}
	for _, valid := range valids {
		err = validateIp("MY_VAR", valid)
		assert.Nil(t, err, "expected %s to be valid", valid)
	}
}

func TestValidatorPort(t *testing.T) {
	err := validatePort("MY_VAR", "ABC")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=ABC is not a valid port: 0-65535", err.Error())
	}

	err = validatePort("MY_VAR", "65535")
	assert.Nil(t, err)
}

func TestValidatorHostport(t *testing.T) {
	err := validateHostport("MY_VAR", "test")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=test is not a valid hostport: address test: missing port in address", err.Error())
	}
	err = validateHostport("MY_VAR", "test:bb")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=test:bb is not a valid hostport: invalid port (bb): 0-65535", err.Error())
	}

	err = validateHostport("MY_VAR", "192.168.0.1:100")
	assert.Nil(t, err)
}

func TestValidatorCidr(t *testing.T) {
	err := validateCidr("MY_VAR", "test")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=test is not a valid CIDR: invalid CIDR address: test", err.Error())
	}

	err = validateCidr("MY_VAR", "192.168.0.1/24")
	assert.Nil(t, err)

	err = validateCidr("MY_VAR", "::/0")
	assert.Nil(t, err)
}

func TestValidateFile(t *testing.T) {
	err := validateFile("MY_VAR", ".fakeFile")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=.fakeFile is not a valid file: CreateFile .fakeFile: The system cannot find the file specified.", err.Error())
	}
	err = validateFile("MY_VAR", "testdata/")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=testdata/ exists but is not a file", err.Error())
	}

	err = validateFile("MY_VAR", ".env")
	assert.Nil(t, err)
}

func TestValidateDir(t *testing.T) {
	err := validateDir("MY_VAR", ".fakeFile")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=.fakeFile is not a valid directory: CreateFile .fakeFile: The system cannot find the file specified.", err.Error())
	}
	err = validateDir("MY_VAR", ".env")
	if assert.Error(t, err) {
		assert.Equal(t, "MY_VAR=.env exists but is not a directory", err.Error())
	}

	err = validateDir("MY_VAR", "testdata/")
	assert.Nil(t, err)
}
