package logger

// using deffrent method's for see log messages on,
// deffrent labels

func (l *Logger) Info(prefix, content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[INFO] (%s) %s", prefix, content)
}

func (l *Logger) Error(prefix, content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[ERROR] (%s) %s", prefix, content)
}

func (l *Logger) Warn(prefix, content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[WARNING] (%s) %s", prefix, content)
}
