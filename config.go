package confik

import "reflect"

type Config struct {
	UseEnvFile       bool
	EnvFilePath      string
	EnvFileOverride  bool
	CustomValidators map[string]FieldValidator
	CustomConverters map[reflect.Type]TypeConverter
}

var DefaultConfig = Config{
	UseEnvFile:      true,
	EnvFilePath:     "",
	EnvFileOverride: false,
}
