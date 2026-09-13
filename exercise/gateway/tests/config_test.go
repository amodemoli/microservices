package tests

import (
	"os"
	"testing"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
)

// TestConfig function, for testing config loader
func TestConfig(t *testing.T) {
	// config example file path (fake-path)
	cnfPath := "./config_test.jsonc"
	jsonc := `{
    // add yours...
    "services": {
        // services host http address
        "user": "fake:5051",
        "auth": "http://fake:5050"
    }
}`

	// create config file if not exists
	if err := os.WriteFile(cnfPath, []byte(jsonc), 0644); err != nil {
		// cannot create, show as error
		t.Errorf("cannot write config file, err: %v", err)
		return // for skip other section's on testing and get result
	}

	// make a fake application for ignore config loader needed app
	var fakeApp *fastic.App

	err, cnf := config.Load(fakeApp, cnfPath)
	if err != nil {
		t.Errorf("cannot load config: %v", err)
		return // for skip other section's of testing
	}

	if cnf.Services["auth"] != "http://fake:5050" {
		t.Errorf("wrong result! config: %s", cnf.Services["auth"])
	}

	if cnf.Services["user"] != "fake:5051" {
		t.Errorf("wrong result! config: %s", cnf.Services["user"])
	}
}
