package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/amodemoli/fastic/core/color"
	"github.com/amodemoli/fastic/core/fastic"
	"github.com/amodemoli/microservices/exercise/gateway/internal/config"
	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
)

type Logger struct {
	writer       *bufio.Writer
	logger       *log.Logger
	cnf          *config.Config
	app          *fastic.App
	file         *os.File
	mu           sync.Mutex
	stopCh       chan struct{}
	wg           sync.WaitGroup
	canAutoFlush bool
	canBackup    bool
}

func New(app *fastic.App, cnf *config.Config, customPath ...string) *Logger {

	path := getPath(cnf, customPath...)

	file, err := readyFile(path)
	if err != nil {
		if app.Env.DevelopemtMode {
			helpers.Print(app, color.Red, "ERROR", err.Error())
			return &Logger{}
		}
		log.Fatalf("error: %v", err)
	}

	writer := bufio.NewWriterSize(file, 64*1024)
	logger := log.New(writer, "", log.Ltime|log.Ldate)

	autoFlush, ok := cnf.General["auto_flush"].(bool)
	if !ok {
		autoFlush = true
	}

	backup, ok := cnf.General["enable_loging_backup"].(bool)
	if !ok {
		backup = true
	}

	lg := &Logger{
		writer:       writer,
		logger:       logger,
		cnf:          cnf,
		app:          app,
		file:         file,
		stopCh:       make(chan struct{}),
		mu:           sync.Mutex{},
		canAutoFlush: autoFlush,
		canBackup:    backup,
	}

	if lg.canAutoFlush {
		lg.wg.Add(1)
		go lg.autoFlush()
	}

	if lg.canBackup {
		lg.startBackup()
	}

	return lg
}

func getPath(cnf *config.Config, customPath ...string) string {
	// make path variable with default velue because if cannot get it log file path from,
	// config file. use this
	path := "../../logs/gateway.log"

	if len(customPath) != 0 && customPath[0] != "" {
		path = customPath[0]
	}

	cnfPath, ok := cnf.General["loging_file"].(string)
	if ok {
		if cnfPath != "" {
			path = cnfPath
		}
	}

	return path
}

func readyFile(path string) (*os.File, error) {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return &os.File{}, fmt.Errorf("%w [logger/new.go:readyFile():MkdirAll()]", err)
	}

	// open file with custom flags
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return &os.File{}, fmt.Errorf("%w [logger/new.go:readyFile():OpenFile()]", err)
	}

	return file, nil
}
