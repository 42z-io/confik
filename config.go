package confik

import "reflect"

type Config[T any] struct {
	UseEnvFile       bool
	EnvFilePath      string
	EnvFileOverride  bool
	CustomValidators map[string]FieldValidator
	CustomConverters map[reflect.Type]TypeConverter
	DefaultValue     *T
}

func DefaultConfig[T any]() Config[T] {
	return Config[T]{
		UseEnvFile:      true,
		EnvFilePath:     "",
		EnvFileOverride: false,
		DefaultValue:    nil,
	}
}
