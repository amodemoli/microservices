package tests

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers/logger"
)

func TestLogger(t *testing.T) {

	path := "./loging_test.log"

	fakeCnf := &config.Config{
		General: map[string]any{
			"loging_file":                     path,
			"auto_flush":                      true,
			"development_logger_flush_ticker": "2s",
			"production_logger_flush_ticker":  "2s",
		},
	}

	fakeApp := &fastic.App{
		Env: &fastic.Env{
			DevelopemtMode: true,
		},
	}

	err := testingLogger(fakeApp, fakeCnf, path, 1)
	if err != nil {
		t.Error(err)
	}

	err = testingLogger(fakeApp, fakeCnf, path, 2)
	if err != nil {
		t.Error(err)
	}

}

func testingLogger(app *fastic.App, cnf *config.Config, path string, step int) error {

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("cannot remove: %v", err)
	}

	lg := logger.New(app, cnf)

	switch step {
	case 1:
		lg.Info("testing with flush")
		// flush now
		lg.Flush()

		err, isEmpty := isEmpty(path)
		if err != nil {
			return fmt.Errorf("cannot read: %v", err)
		}

		if isEmpty {
			return errors.New("content is empty")
		}

	case 2:
		lg.Info("testing without flush")
		// testing autoflush (dont flush with lg.Flush)

		// sleeped 3 seconds (auto-flush delay is two seconds but im sleeping for three seconds)
		time.Sleep(3 * time.Second)

		err, isEmpty := isEmpty(path)
		if err != nil {
			return fmt.Errorf("cannot read: %v", err)
		}

		if isEmpty {
			return errors.New("content is empty")
		}

	default:
		return fmt.Errorf("%v is invalid step, only accept's 1||2", step)
	}

	return nil
}

func isEmpty(path string) (error, bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return err, true
	}

	var isEmpty bool

	if string(content) == "" {
		isEmpty = true
	}

	return nil, isEmpty
}
