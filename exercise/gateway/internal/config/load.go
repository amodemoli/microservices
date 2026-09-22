package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amodemoli/fastic/core/fastic"
)

// Config model, for save it jsonc config result here,
// used maps because map dont need to edit result's here only can edit key on json,
// and key && value are resets here dont need edit here before edit keys on json file
type Config struct {
	// Registered Services url and name are saves here
	// used string for value type because i know only string can save it as value!
	// value is service http server url.
	Services map[string]string `json:"services"`

	// Middleware's general settings are saves here,
	// setting's will be used only in middlewares
	Middewares map[string]any `json:"middlewares"`

	// General settings of proxy (gateway),
	// e.g: timeout of pinging services on health path
	General map[string]any `json:"general"`
}

// Load function maked for return error and config model,
// on returned config model save's jsonc config values
// [*] you need add path of config.jsonc file and example of app model to function for work
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
