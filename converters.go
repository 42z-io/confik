package confik

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func convertError(envName string, value string, kind reflect.Kind, err error) error {
	return fmt.Errorf("%s=%s is not a valid %s: %w", envName, value, kind, err)
}

type TypeConverter = func(fc *FieldConfig, fieldValue string, rv reflect.Value) error

func handleUint8(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 8)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint8, err)
	}
	rv.SetUint(i)
	return nil
}

func handleUint16(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 16)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint16, err)
	}
	rv.SetUint(i)
	return nil
}
func handleUint32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint32, err)
	}
	rv.SetUint(i)
	return nil
}
func handleUint64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseUint(fieldValue, 10, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Uint64, err)
	}
	rv.SetUint(i)
	return nil
}

func handleInt8(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 8)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int8, err)
	}
	rv.SetInt(i)
	return nil
}

func handleInt16(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 16)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int16, err)
	}
	rv.SetInt(i)
	return nil
}

func handleInt32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int32, err)
	}
	rv.SetInt(i)
	return nil
}

func handleInt64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseInt(fieldValue, 10, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Int64, err)
	}
	rv.SetInt(i)
	return nil
}

func handleFloat32(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseFloat(fieldValue, 32)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Float32, err)
	}
	rv.SetFloat(i)
	return nil
}

func handleFloat64(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	i, err := strconv.ParseFloat(fieldValue, 64)
	if err != nil {
		return convertError(fc.Name, fieldValue, reflect.Float64, err)
	}
	rv.SetFloat(i)
	return nil
}

func handleBool(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
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

func handleString(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	rv.SetString(fieldValue)
	return nil
}

func handleStringSlice(fc *FieldConfig, fieldValue string, rv reflect.Value) error {
	if rv.Type() != reflect.SliceOf(reflect.TypeOf(fieldValue)) {
		return fmt.Errorf("%s is invalid: only string slices are supported", fc.Name)
	}
	slice := strings.Split(fieldValue, ",")
	rv.Set(reflect.ValueOf(slice))
	return nil
}

var typeConverters = map[reflect.Kind]TypeConverter{
	reflect.Uint8:   handleUint8,
	reflect.Uint16:  handleUint16,
	reflect.Uint32:  handleUint32,
	reflect.Uint64:  handleUint64,
	reflect.Int8:    handleInt8,
	reflect.Int16:   handleInt16,
	reflect.Int32:   handleInt32,
	reflect.Int64:   handleInt64,
	reflect.Float32: handleFloat32,
	reflect.Float64: handleFloat64,
	reflect.Bool:    handleBool,
	reflect.String:  handleString,
	reflect.Slice:   handleStringSlice,
}
