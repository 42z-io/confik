package confik

import "reflect"

type Config[T any] struct {
	UseEnvFile      bool
	EnvFilePath     string
	EnvFileOverride bool
	Validators      map[string]Validator
	Parsers         map[reflect.Type]Parser
	DefaultValue    *T
}

func DefaultConfig[T any]() Config[T] {
	return Config[T]{
		UseEnvFile:      true,
		EnvFilePath:     "",
		EnvFileOverride: false,
		DefaultValue:    nil,
	}
}
