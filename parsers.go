package confik

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Parser is the type a function must implement to provide type conversion for types.
type Parser = func(fc *FieldConfig, fieldValue string, rv reflect.Value) error

func convertError(envName string, value string, kind reflect.Kind, err error) error {
	return fmt.Errorf("%s=%s is not a valid %s: %w", envName, value, kind, err)
}

func parseUint(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, strconv.IntSize)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint, err)
	}
	rv.SetUint(i)
	return nil
}

func parseUint8(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 8)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint8, err)
	}
	rv.SetUint(i)
	return nil
}

func parseUint16(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 16)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint16, err)
	}
	rv.SetUint(i)
	return nil
}
func parseUint32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint32, err)
	}
	rv.SetUint(i)
	return nil
}
func parseUint64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint64, err)
	}
	rv.SetUint(i)
	return nil
}

func parseInt(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, strconv.IntSize)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int, err)
	}
	rv.SetInt(i)
	return nil
}

func parseInt8(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 8)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int8, err)
	}
	rv.SetInt(i)
	return nil
}

func parseInt16(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 16)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int16, err)
	}
	rv.SetInt(i)
	return nil
}

func parseInt32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int32, err)
	}
	rv.SetInt(i)
	return nil
}

func parseInt64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int64, err)
	}
	rv.SetInt(i)
	return nil
}

func parseFloat32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseFloat(fieldValue, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Float32, err)
	}
	rv.SetFloat(i)
	return nil
}

func parseFloat64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Float64, err)
	}
	rv.SetFloat(i)
	return nil
}

func parseBool(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	lower := strings.ToLower(fieldValue)
	if lower == "0" || lower == "false" || lower == "no" {
		rv.SetBool(false)
	} else if lower == "1" || lower == "true" || lower == "yes" {
		rv.SetBool(true)
	} else {
		return fmt.Errorf("%s=%s is not a valid bool: expects yes/no, true/false, 0/1", fc.Name, fieldValue)
	}
	return nil
}

func parseString(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	rv.SetString(fieldValue)
	return nil
}

func handleSlice(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	strSlice := strings.Split(fieldValue, ",")
	var unsliced = rv.Type().Elem()
	converter, exists := kindParsers[unsliced.Kind()]
	if !exists {
		return fmt.Errorf("%s is invalid: %s is not supported", fc.Name, rv.Type())
	}
	var data = reflect.MakeSlice(rv.Type(), 0, len(strSlice))
	for _, v := range strSlice {
		rv2 := reflect.New(unsliced).Elem()
		err := converter(fc, v, rv2)
		if err != nil {
			return fmt.Errorf("%s=%s is not a valid %s: %w", fc.Name, fieldValue, rv.Type(), errors.Unwrap(err))
		}
		data = reflect.Append(data, rv2)
	}
	rv.Set(data)
	return nil
}

func parseUrl(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	u, err := url.ParseRequestURI(fieldValue)
	if err != nil {
		return fmt.Errorf("%s=%s invalid url.URL: %w", fc.Name, fieldValue, err)
	}
	rv.Set(reflect.ValueOf(*u))
	return nil
}

func parseTime(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	t, err := time.Parse(time.RFC3339, fieldValue)
	if err != nil {
		return fmt.Errorf("%s=%s invalid time.Time: %w", fc.Name, fieldValue, err)
	}
	rv.Set(reflect.ValueOf(t))
	return nil
}

func parseDuration(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	t, err := time.ParseDuration(fieldValue)
	if err != nil {
		return fmt.Errorf("%s=%s invalid time.Duration: %w", fc.Name, fieldValue, err)
	}
	rv.Set(reflect.ValueOf(t))
	return nil
}

var typeParsers = map[reflect.Type]Parser{
	reflect.TypeOf((*url.URL)(nil)).Elem():       parseUrl,
	reflect.TypeOf((*time.Time)(nil)).Elem():     parseTime,
	reflect.TypeOf((*time.Duration)(nil)).Elem(): parseDuration,
}

var kindParsers = map[reflect.Kind]Parser{
	reflect.Uint:    parseUint,
	reflect.Uint8:   parseUint8,
	reflect.Uint16:  parseUint16,
	reflect.Uint32:  parseUint32,
	reflect.Uint64:  parseUint64,
	reflect.Int:     parseInt,
	reflect.Int8:    parseInt8,
	reflect.Int16:   parseInt16,
	reflect.Int32:   parseInt32,
	reflect.Int64:   parseInt64,
	reflect.Float32: parseFloat32,
	reflect.Float64: parseFloat64,
	reflect.Bool:    parseBool,
	reflect.String:  parseString,
}
