package confik

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyCustomType struct {
	Value string
}

func parseMyCustomType(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	var res = MyCustomType{
		Value: fieldValue,
	}
	rv.Set(reflect.ValueOf(res))
	return nil
}

type CustomConverterType struct {
	MyField MyCustomType
}

func TestLoadFromEnvCustomConverter(t *testing.T) {
	os.Clearenv()
	os.Setenv("MY_FIELD", "hello")
	custom, err := LoadFromEnv[CustomConverterType](Config[CustomConverterType]{
		UseEnvFile: false,
		Parsers: map[reflect.Type]Parser{
			reflect.TypeOf(MyCustomType{}): parseMyCustomType,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, custom.MyField.Value, "hello")
}

func TestUint8(t *testing.T) {
	var res uint8
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseUint8(fc, "-100", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-100 is not a valid uint8: strconv.ParseUint: parsing \"-100\": invalid syntax", err.Error())
	}
	err = parseUint8(fc, "100", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, uint8(100))
}

func TestUint16(t *testing.T) {
	var res uint16
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseUint16(fc, "-100", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-100 is not a valid uint16: strconv.ParseUint: parsing \"-100\": invalid syntax", err.Error())
	}
	err = parseUint16(fc, "3233", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, uint16(3233))
}

func TestUint32(t *testing.T) {
	var res uint32
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseUint32(fc, "-100", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-100 is not a valid uint32: strconv.ParseUint: parsing \"-100\": invalid syntax", err.Error())
	}
	err = parseUint32(fc, "66444", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, uint32(66444))
}

func TestUint64(t *testing.T) {
	var res uint64
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseUint64(fc, "-100", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-100 is not a valid uint64: strconv.ParseUint: parsing \"-100\": invalid syntax", err.Error())
	}
	err = parseUint64(fc, "18446744073709551615", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, uint64(18446744073709551615))
}

func TestInt8(t *testing.T) {
	var res int8
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseInt8(fc, "-255", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-255 is not a valid int8: strconv.ParseInt: parsing \"-255\": value out of range", err.Error())
	}
	err = parseInt8(fc, "100", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, int8(100))
}

func TestInt16(t *testing.T) {
	var res int16
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseInt16(fc, "-32769", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=-32769 is not a valid int16: strconv.ParseInt: parsing \"-32769\": value out of range", err.Error())
	}
	err = parseInt16(fc, "-32768", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, int16(-32768))
}

func TestInt32(t *testing.T) {
	var res int32
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseInt32(fc, "2147483648", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=2147483648 is not a valid int32: strconv.ParseInt: parsing \"2147483648\": value out of range", err.Error())
	}
	err = parseInt32(fc, "2000000000", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, int32(2000000000))
}

func TestInt64(t *testing.T) {
	var res int64
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseInt64(fc, "18446744073709551615", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=18446744073709551615 is not a valid int64: strconv.ParseInt: parsing \"18446744073709551615\": value out of range", err.Error())
	}
	err = parseInt64(fc, "-9223372036854775808", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, int64(-9223372036854775808))
}

func TestFloat32(t *testing.T) {
	var res float32
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseFloat32(fc, "abc", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=abc is not a valid float32: strconv.ParseFloat: parsing \"abc\": invalid syntax", err.Error())
	}
	err = parseFloat32(fc, "3.1459", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, float32(3.1459))
}

func TestFloat64(t *testing.T) {
	var res float64
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseFloat64(fc, "1.7e+309", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=1.7e+309 is not a valid float64: strconv.ParseFloat: parsing \"1.7e+309\": value out of range", err.Error())
	}
	err = parseFloat64(fc, "1.7e+308", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, float64(1.7e+308))
}

func TestBool(t *testing.T) {
	var res bool
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseBool(fc, "g", rv)
	if assert.Error(t, err) {
		assert.Equal(t, "test=g is not a valid bool: expects yes/no, true/false, 0/1", err.Error())
	}
	err = parseBool(fc, "true", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, true)
	err = parseBool(fc, "1", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, true)
	err = parseBool(fc, "yes", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, true)
	err = parseBool(fc, "false", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, false)
	err = parseBool(fc, "no", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, false)
	err = parseBool(fc, "0", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, false)
}

func TestString(t *testing.T) {
	var res string
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := parseString(fc, "string", rv)
	assert.Nil(t, err)
	assert.Equal(t, res, "string")
}

func TestStringSlice(t *testing.T) {
	var res []string
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := handleSlice(fc, "string1,string2", rv)
	assert.Nil(t, err)
	assert.Equal(t, []string{"string1", "string2"}, res)
}

func TestIntSlice(t *testing.T) {
	var res []int
	rv := reflect.ValueOf(&res).Elem()
	fc := &FieldConfig{
		ConfigTag: NewConfigTag("test"),
		Validate:  nil,
	}
	err := handleSlice(fc, "1,2", rv)
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2}, res)
}
