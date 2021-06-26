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

type Config struct {
	ConfigFile string
	Debug      bool   `toml:"debug" yaml:"debug" json:"debug"`
	BatchMode  bool   `toml:"batch" yaml:"batch" json:"batch"`
	DataFile   string `toml:"data_file" yaml:"data_file" json:"data_file"`
}

// loading configuration parameters from a file
func (c Config) loadConfFile(path string) error {
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
	flag.Parse()

	if result.ConfigFile != "" {
		if err := result.loadConfFile(result.ConfigFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	flag.BoolVar(&result.Debug, "debug", GetEnvBool("DEBUG", false), "Output of detailed debugging information")
	flag.BoolVar(&result.BatchMode, "batch", GetEnvBool("BATCH", false), "Run in batch/no interactive mode")
	flag.StringVar(&result.DataFile, "data", GetEnv("DATA", ""), "The path to the file *.scv with data (recuired)")
	// flag.StringVar(&result.Hash, "hash", GetEnv("HASH", "md5"), "HASH method: [md5, crc32], default md5")
	// flag.Uint64Var(&result.MaxSize, "max-size", GetEnvUInt("MAX_SIZE", 0), "limit maximum file size for checking, 0 - disable")
	// flag.Uint64Var(&result.DFactor, "double-factor", GetEnvUInt("D_FACTOR", 1), "double factor, default > 1")
	// flag.BoolVar(&result.Size, "size", GetEnvBool("SIZE", true), "Compare files by size")
	flag.Parse()

	// Determine the initial directories.
	// result.Dirs = flag.Args()
	// if len(result.Dirs) == 0 {
	// 	result.Dirs = []string{"."}
	// }

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
