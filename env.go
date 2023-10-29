package confik

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadEnvFile(cfg Config) (map[string]string, error) {
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

	// No .env found - return empty
	if envPath == "" {
		envMap := make(map[string]string)
		return envMap, nil
	}

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

	file, err := os.Open(envPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	kv, err := ParseEnvFile(scanner)
	if err != nil {
		return nil, err
	}
	for k, v := range kv {
		_, exists := os.LookupEnv(k)
		if cfg.EnvFileOverride || !exists {
			os.Setenv(k, v)
		}
	}
	return kv, nil
}

func FindEnvFile() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var lastPath = filepath.Clean(path)
	for {
		var checkPath = filepath.Join(path, ".env")
		stat, err := os.Stat(checkPath)
		if os.IsNotExist(err) {
			path = filepath.Dir(path)
			if path == lastPath {
				return "", nil
			}
			lastPath = path
		} else if stat.IsDir() {
			return "", fmt.Errorf("environment file is a directory: %s", checkPath)
		} else if err != nil {
			return "", err
		} else {
			return checkPath, nil
		}
	}
}

func ParseEnvFile(scanner *bufio.Scanner) (map[string]string, error) {
	kv := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		// ignore comments and blank lines
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.Trim(line, " ") == "" {
			continue
		}
		// ensure is in a variable format (VARIABLE=value)
		if !strings.Contains(line, "=") {
			return nil, fmt.Errorf("invalid line in env file: %s", line)
		}
		var key, value string
		parts := strings.SplitN(line, "=", 2)
		key = parts[0]
		if len(parts) == 1 {
			value = ""
		} else {
			unquoted, err := strconv.Unquote(parts[1])
			if err != nil {
				value = parts[1]
			} else {
				value = unquoted
			}
		}
		kv[key] = strings.Trim(value, " ")
	}
	return kv, nil
}
