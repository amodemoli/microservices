package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amodemoli/fastic/core/fastic"
)

type Config struct {
	Services map[string]string `json:"services"`

	PingingTimeout time.Duration `json:"pinging_timeout"`
}

func Load(app *fastic.App, path string) (error, *Config) {
	// validate config file path
	// 1) check length only
	if path == "" {
		// return error and empty config file
		return errors.New("file path cannot be empty"), &Config{}
	}

	// 2) check the file extension
	// get extension and save it to extension variable
	extension := filepath.Ext(path)
	// check extension is valid or not
	if extension != ".jsonc" {
		// extension is invalid! return error and exit
		return fmt.Errorf("%s file extension is not supported, use .jsonc", extension), &Config{}
	}

	// read data of path file
	data, err := os.ReadFile(path)
	// check error exists or not
	if err != nil {
		// return error and exit from func
		// only return raw error dont need to make new error message
		return err, &Config{}
	}

	// remove comments (start's with //) on json file (because maby user used jsonc)
	filtered := filterJson(string(data))

	// create variable for save config
	// result's there:
	var cnf Config

	// unmarshal config json content to struct
	if err := json.Unmarshal([]byte(filtered), &cnf); err != nil {
		// cannot unmarshal, return error and exit from func
		return fmt.Errorf("cannot unmarshal json: %v", err), &Config{}
	}

	// exit from function and return config struct
	// this config struct has json data. =D
	return nil, &cnf
}
