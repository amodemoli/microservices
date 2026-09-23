package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
)

func (l *Logger) startBackup() {

	prefix := "backup-system"
	path, dir, d, maximumSize := l.getFinalVariables()

	err := os.MkdirAll(path, 0755)
	if err != nil {
		// only i can print error on log file
		l.Error(prefix, fmt.Sprintf("cannot create directory %s: %v", path, err))
	}

	newName := filepath.Join(dir, fmt.Sprintf("log.%s.backup", time.Now().Format("2006-01-02-15:04")))

	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for range ticker.C {

			size, err := l.getFileSize(path)
			if err != nil {
				l.Error(prefix, fmt.Sprintf("cannot read size %s: %v", path, err))
			}

			if size < maximumSize {
				continue
			}

			if err := os.Rename(path, newName); err != nil {
				l.Error(prefix, fmt.Sprintf("cannot rename %s: %v", path, err))
				continue
			}

			f, err := os.Create(path)
			if err != nil {
				// dont use return because logger create's automatic new log file after flush request
				// if logger cannot create log file. print's error on terminal
				l.Error(prefix, fmt.Sprintf("cannot create %s: %v", path, err))
				continue
			}

			if err := f.Close(); err != nil {
				l.Error(prefix, fmt.Sprintf("cannot close %s: %v", path, err))
			}
		}
	}()

}

func (l *Logger) getFinalVariables() (path, dir string, ticker time.Duration, maximumSize int64) {
	path = "../../logs/gateway.log"
	dir = "../../logs/backups/"
	ticker = 30 * time.Second
	// maximumSize is byte (500_000 = 500kb || 0.5mg)
	// im used "_" this section ignore's with compiler im used this on my number for spell number section's and ready easy number
	maximumSize = 500_000

	cnfPath, ok := l.cnf.General["loging_file"].(string)
	if ok {
		if cnfPath != "" {
			path = cnfPath
		}
	}

	cnfDir, ok := l.cnf.General["loging_backup_dir"].(string)
	if ok {
		if dir != "" {
			dir = cnfDir
		}
	}

	cnfTicker, err := helpers.ToDuration(l.cnf.General["auto_check_backup"])
	if err == nil {
		// cannot lower than 1 second
		if cnfTicker >= 1*time.Second {
			ticker = cnfTicker
		}
	}

	cnfMaximumSize, ok := l.cnf.General["maximum_log_size"].(float64)
	if ok {
		// cannot lower than 5,000 byte (5kb)
		if cnfMaximumSize >= 5_000 {
			maximumSize = int64(cnfMaximumSize)
		}
	}

	return path, dir, ticker, maximumSize
}

func (l *Logger) getFileSize(path string) (int64, error) {
	f, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return f.Size(), err
}
