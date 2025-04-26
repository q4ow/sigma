package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Version   string
	Verbose   bool
	ConfigDir string
	HomeDir   string
}

func DefaultConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		Version:   "0.1.0",
		Verbose:   false,
		ConfigDir: filepath.Join(home, ".sigma"),
		HomeDir:   home,
	}, nil
}

var GlobalConfig *Config

func Initialize() error {
	cfg, err := DefaultConfig()
	if err != nil {
		return err
	}
	GlobalConfig = cfg
	return nil
}

func SetVerbose(verbose bool) {
	if GlobalConfig != nil {
		GlobalConfig.Verbose = verbose
	}
}

func IsVerbose() bool {
	if GlobalConfig != nil {
		return GlobalConfig.Verbose
	}
	return false
}

func GetVersion() string {
	if GlobalConfig != nil {
		return GlobalConfig.Version
	}
	return "unknown"
}
