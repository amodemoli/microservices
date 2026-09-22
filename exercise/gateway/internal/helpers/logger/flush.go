package logger

import (
	"time"

	"github.com/amodemoli/microservices/exercise/gateway/internal/helpers"
)

func (l *Logger) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.writer != nil {
		l.writer.Flush()
	}

	if l.file != nil {
		l.file.Sync()
	}
}

// Close method maked for close opened log file,
// dont worry this method flush content's to writer after close
func (l *Logger) Close() {
	if l.stopCh != nil {
		close(l.stopCh)
	}

	l.wg.Wait()

	l.mu.Lock()
	defer l.mu.Unlock()

	// changing writer to nil
	if l.writer != nil {
		l.writer.Flush()
	}

	if l.file != nil {
		l.file.Close()
	}
}

// autoFlush is private method because need it run only one time,
// dont need to run any time and developer dont need to access on this section
// Dont forget: custom interval cannot lower than 1 second
func (l *Logger) autoFlush(d ...time.Duration) {
	defer l.wg.Done()

	// get interval from custom fucntion
	interval := l.getInterval(d...)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.Flush()
		case <-l.stopCh:
			// final flush before closing
			l.Flush()
			return
		}
	}
}

func (l *Logger) getInterval(d ...time.Duration) time.Duration {
	var interval time.Duration
	var anyCnfInternal any

	// defferent default values for ticker duration default value,
	// because on development mode developer need to see errors on log very speed. but on production mode developer dont need to speed flush,
	// they need to have +10 second time for flush ticker for performance.
	if l.app.Env.DevelopemtMode {
		interval = 1 * time.Second
		anyCnfInternal = l.cnf.General["development_logger_flush_ticker"]
	} else {
		interval = 25 * time.Second
		anyCnfInternal = l.cnf.General["production_logger_flush_ticker"]
	}

	if len(d) != 0 && d[0] >= 1*time.Second {
		return d[0]
	}

	cnfInternal, err := helpers.ToDuration(anyCnfInternal)
	if err == nil {
		interval = cnfInternal
	}

	return interval
}
