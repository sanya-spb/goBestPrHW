package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/komkom/toml"
	"gopkg.in/yaml.v3"
)

// configuration struct
type Config struct {
	ConfigFile    string
	Debug         bool   `toml:"debug" yaml:"debug" json:"debug"`
	BatchMode     bool   `toml:"batch" yaml:"batch" json:"batch"`
	DataFile      string `toml:"data_file" yaml:"data_file" json:"data_file"`
	FilterTimeout uint64 `toml:"filter_timeout" yaml:"filter_timeout" json:"filter_timeout"`
	LogAccess     string `toml:"log_access" yaml:"log_access" json:"log_access"`
	LogErrors     string `toml:"log_errors" yaml:"log_errors" json:"log_errors"`
}

// loading configuration parameters from a file
func (c *Config) loadConfFile(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if f, err := os.Open(path); err != nil {
			return err
		} else {
			defer f.Close()

			switch ext := strings.ToLower(filepath.Ext(path)); ext {
			case ".toml":
				if err := json.NewDecoder(toml.New(f)).Decode(&c); err != nil {
					return err
				}
				// toml.Unmarshal([]byte(""), c)
			case ".yaml", ".yml":
				if err := yaml.NewDecoder(f).Decode(&c); err != nil {
					return err
				}
			case ".json":
				if err := json.NewDecoder(f).Decode(&c); err != nil {
					return err
				}
			default:
				return fmt.Errorf("Unknown file format: %s", path)
			}
		}
	} else {
		return fmt.Errorf("The configuration file is missing: %s", path)
	}
	return nil
}

// Init config
func NewConfig() *Config {
	var result *Config = new(Config)
	flag.StringVar(&result.ConfigFile, "config", GetEnv("CONFIG", ""), "Configuration settings file")
	flag.BoolVar(&result.Debug, "debug", GetEnvBool("DEBUG", false), "Output of detailed debugging information")
	flag.BoolVar(&result.BatchMode, "batch", GetEnvBool("BATCH", false), "Run in batch/no interactive mode")
	flag.StringVar(&result.DataFile, "data", GetEnv("DATA", ""), "The path to the file *.scv with data (recuired)")
	flag.StringVar(&result.LogAccess, "log-access", GetEnv("LOG_ACCESS", "./access.log"), "Log file for user input")
	flag.StringVar(&result.LogErrors, "log-errors", GetEnv("LOG_ERRORS", "./errors.log"), "Log file for errors")
	flag.Uint64Var(&result.FilterTimeout, "filter-timeout", GetEnvUInt("FILTER_TIMEOUT", 1000), "Timeout to filtering data, ms")
	flag.Parse()

	if result.ConfigFile != "" {
		if err := result.loadConfFile(result.ConfigFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	return result
}

// Get uint value from ENV
func GetEnvUInt(key string, defaultVal uint64) uint64 {
	if envVal, ok := os.LookupEnv(key); ok {
		if envBool, err := strconv.ParseUint(envVal, 10, 64); err == nil {
			return envBool
		}
	}
	return defaultVal
}

// Get string value from ENV
func GetEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

// Get bool value from ENV
func GetEnvBool(key string, defaultVal bool) bool {
	if envVal, ok := os.LookupEnv(key); ok {
		if envBool, err := strconv.ParseBool(envVal); err == nil {
			return envBool
		}
	}
	return defaultVal
}
