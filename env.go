package confik

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadEnvFile[T any](cfg Config[T]) (map[string]string, error) {
	// find the .env file
	//
	var envPath string
	if cfg.EnvFilePath == "" {
		foundPath, err := FindEnvFile()
		if err != nil {
			return nil, err
		}
		envPath = foundPath
	} else {
		envPath = cfg.EnvFilePath
	}

	// No .env found or provided - return empty map
	if envPath == "" {
		envMap := make(map[string]string)
		return envMap, nil
	}

	// Check if the .env file exists and is not a directory
	stat, err := os.Stat(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("environment file does not exist: %s", envPath)
		}
		return nil, err
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("environment file is a directory: %s", envPath)
	}

	// Open and parse the env file
	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	kv, err := ParseEnvFile(bufio.NewScanner(file))
	if err != nil {
		return nil, err
	}

	// Add the discovered environment variables in the .env to the environment
	for k, v := range kv {
		_, exists := os.LookupEnv(k)
		if cfg.EnvFileOverride || !exists {
			os.Setenv(k, v)
		}
	}
	return kv, nil
}

func FindEnvFile() (string, error) {
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

func parseEnvironmentVariableExpression(expression string) (string, string, error) {
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

func ParseEnvFile(scanner *bufio.Scanner) (map[string]string, error) {
	// map to store all the key => value pairs we find in the ".env" file
	kv := make(map[string]string)
	for scanner.Scan() {
		expression := scanner.Text()
		// ignore comments and blank lines
		if strings.HasPrefix(expression, "#") || strings.HasPrefix(expression, "//") || strings.TrimSpace(expression) == "" {
			continue
		}

		// parse the environment variable
		key, value, err := parseEnvironmentVariableExpression(expression)
		if err != nil {
			return nil, err
		}
		// update the store
		kv[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return kv, nil
}
