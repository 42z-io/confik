package confik

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetting(t *testing.T) {
	_, _, err := parseSetting("hello")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid setting hello: invalid syntax", err.Error())
	}

	key, value, err := parseSetting("hello=world")
	assert.Nil(t, err)
	assert.Equal(t, "hello", key)
	assert.Equal(t, "world", value)
}

func TestParseTag(t *testing.T) {
	_, err := parseTag("")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid tag: empty tag", err.Error())
	}
	_, err = parseTag("hello,flag1,flag2,option==")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid tag: invalid setting option==: invalid syntax", err.Error())
	}
	tag, err := parseTag("name,flag1,flag2,option1=opt1,option2=opt2")
	assert.Nil(t, err)
	assert.Equal(t, "name", tag.Name)
	assert.Equal(t, []string{"flag1", "flag2"}, tag.Flags)
	assert.Equal(t, map[string]string{"option1": "opt1", "option2": "opt2"}, tag.Settings)
}

func TestParseEnvTag(t *testing.T) {
	_, err := parseEnvTag("")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid env tag: invalid tag: empty tag", err.Error())
	}

	_, err = parseEnvTag("@@")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid env tag: invalid environment variable name: @@ must be [A-Z0-9_]+", err.Error())
	}

	_, err = parseEnvTag("NAME,flag1")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid env tag: unknown flag flag1", err.Error())
	}

	_, err = parseEnvTag("NAME,option1=opt1")
	if assert.Error(t, err) {
		assert.Equal(t, "invalid env tag: unknown setting option1", err.Error())
	}

	tag, err := parseEnvTag("NAME,optional,unset,default=DEFAULT,validate=validator")
	assert.Nil(t, err)
	assert.Equal(t, "NAME", tag.Name)
	assert.Equal(t, true, tag.Optional)
	assert.Equal(t, true, tag.Unset)
	assert.Equal(t, "DEFAULT", *tag.Default)
	assert.Equal(t, "validator", *tag.Validator)
}
