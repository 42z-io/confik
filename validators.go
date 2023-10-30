package confik

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
)

type Validator = func(envName, value string) error

func validateUri(envName string, value string) error {
	if _, err := url.ParseRequestURI(value); err != nil {
		return fmt.Errorf("%s=%s is not a URI: %w", envName, value, err)
	}
	return nil
}

func validateIp(envName string, value string) error {
	if ip := net.ParseIP(value); ip == nil {
		return fmt.Errorf("%s=%s is not a valid IP: invalid format", envName, value)
	}
	return nil
}

func validatePort(envName string, value string) error {
	_, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid port: 0-65535", envName, value)
	}
	return nil
}

func validateHostport(envName string, value string) error {
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid hostport: %w", envName, value, err)
	}
	err = validatePort(envName, port)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid hostport: invalid port (%s): 0-65535", envName, value, port)
	}
	return nil
}

func validateCidr(envName string, value string) error {
	_, _, err := net.ParseCIDR(value)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid CIDR: %w", envName, value, err)
	}
	return nil
}

func validateFile(envName string, value string) error {
	stat, err := os.Stat(value)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid file: %w", envName, value, err)
	}
	if stat.IsDir() {
		return fmt.Errorf("%s=%s exists but is not a file", envName, value)
	}
	return nil
}

func validateDir(envName string, value string) error {
	stat, err := os.Stat(value)
	if err != nil {
		return fmt.Errorf("%s=%s is not a valid directory: %w", envName, value, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s=%s exists but is not a directory", envName, value)
	}
	return nil
}

var fieldValidators = map[string]Validator{
	"uri":      validateUri,
	"ip":       validateIp,
	"port":     validatePort,
	"hostport": validateHostport,
	"cidr":     validateCidr,
	"file":     validateFile,
	"dir":      validateDir,
}
