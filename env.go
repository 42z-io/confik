package confik

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// loadEnvFile will locate and load the environment file into a map[string]string
//
// loadEnvFile will update the current environment with the files found in the environment file
func loadEnvFile[T any](cfg Config[T]) (map[string]string, error) {
	var envPath string
	if cfg.EnvFilePath == "" {
		foundPath, err := findEnvFile()
		if err != nil {
			return nil, err
		}
		envPath = foundPath
	} else {
		envPath = cfg.EnvFilePath
	}

	// no .env found or provided - return empty map
	if envPath == "" {
		envMap := make(map[string]string)
		return envMap, nil
	}

	// check if the .env file exists
	stat, err := os.Stat(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("environment file does not exist: %s", envPath)
		}
		return nil, err
	}

	// check if the .env file is a directory
	if stat.IsDir() {
		return nil, fmt.Errorf("environment file is a directory: %s", envPath)
	}

	// open and parse the env file
	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	kv, err := parseEnvFile(file)
	if err != nil {
		return nil, err
	}

	// add the discovered environment variables in the environment file to the environment
	for k, v := range kv {
		_, exists := os.LookupEnv(k)
		if cfg.EnvFileOverride || !exists {
			os.Setenv(k, v)
		}
	}
	return kv, nil
}

// findEnvFile will locate the .env file by looking in the current directory and recursing up the directory structure
func findEnvFile() (string, error) {
	// Start looking for the ".env" in the current directory
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var lastPath = filepath.Clean(path)
	for {
		var checkPath = filepath.Join(path, ".env")
		stat, err := os.Stat(checkPath)
		// If we cant find the ".env" in this directory look in the parent directory
		if os.IsNotExist(err) {
			path = filepath.Dir(path)
			if path == lastPath {
				return "", nil
			}
			lastPath = path
		} else if err != nil {
			return "", err
		} else if stat.IsDir() {
			return "", fmt.Errorf("environment file is a directory: %s", checkPath)
		} else {
			return checkPath, nil
		}
	}
}

// parseEnvVar will parse an environment variable in the format NAME=VALUE.
func parseEnvVar(expression string) (string, string, error) {
	// ensure format of the expression is correct
	if !strings.Contains(expression, "=") {
		return "", "", fmt.Errorf("invalid expression in env file: %s", expression)
	}

	// split the variable NAME=value [NAME, value]
	parts := strings.SplitN(expression, "=", 2)
	variable := strings.Trim(parts[0], " ")

	// remove any quotes
	unquoted, err := strconv.Unquote(parts[1])
	if err != nil {
		return variable, parts[1], nil
	}
	return variable, unquoted, nil
}

// parseEnvFile will convert an environment file into a map[string]string
//
// Expects the file in the format:
//
//	MY_VARIABLE=MY_NAME
//	OTHER_VARIABLE="QUOTED_VALUE"
//
// Notes:
//   - Quoted values will be unquoted
//   - Blank lines will be ingored
//   - Comments (starting with // or #) will be ignored
//   - Whitespace around variables and their values will be stripped
func parseEnvFile(reader io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(reader)
	// map to store all the key => value pairs we find in the ".env" file
	kv := make(map[string]string)
	for scanner.Scan() {
		expression := scanner.Text()
		// ignore comments and blank lines
		if strings.HasPrefix(expression, "#") || strings.HasPrefix(expression, "//") || strings.TrimSpace(expression) == "" {
			continue
		}

		key, value, err := parseEnvVar(expression)
		if err != nil {
			return nil, err
		}

		// update the store
		kv[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return kv, nil
}
