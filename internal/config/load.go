package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func Load() (*Config, error) {
	cfgFile, err := os.Open(FileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrConfigNotFound, FileName)
		}

		return nil, err
	}

	defer func() { _ = cfgFile.Close() }()

	decoder := json.NewDecoder(cfgFile)
	decoder.DisallowUnknownFields()

	var cfg Config

	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", FileName, err)
	}

	if err = validate(&cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", FileName, err)
	}

	return &cfg, nil
}
