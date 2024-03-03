package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func LoadConfig(file string) (*Config, error) {
	var result Config
	configFile, err := os.Open(file)

	if errors.Is(err, os.ErrNotExist) {
		path := file[0:strings.LastIndex(file, "/")]

		err = os.MkdirAll(path, 0770)

		if err != nil {
			return nil, err
		}

		configFile, err = os.Create(file)

		if err != nil {
			return nil, err
		}

		defer configFile.Close()

		var defaultConfig = CreateDefault()

		bytes, err := json.MarshalIndent(defaultConfig, " ", "	")

		if err != nil {
			return nil, err
		}

		_, err = configFile.WriteString(string(bytes))
		if err != nil {
			return nil, err
		}

		return defaultConfig, nil
	} else if err != nil {
		return nil, err
	}

	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
