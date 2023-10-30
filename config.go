package confik

import "reflect"

// Config[T] is the configuration for reading environment variables.
type Config[T any] struct {
	UseEnvFile      bool                    // read from an environment file on disk?
	EnvFilePath     string                  // custom path to the environment file (otherwise search for ".env")
	EnvFileOverride bool                    // should variables found in the env file override environment variables?
	Validators      map[string]Validator    // a map of custom validators to be used by the loader
	Parsers         map[reflect.Type]Parser // a map of custom type parsers to be used by the loader
	DefaultValue    *T                      // default values to use if they do not exist in the environment
}

// DefaultConfig will create a new [Config] with the default values.
func DefaultConfig[T any]() Config[T] {
	return Config[T]{
		UseEnvFile:      true,
		EnvFilePath:     "",
		EnvFileOverride: false,
		DefaultValue:    nil,
	}
}
