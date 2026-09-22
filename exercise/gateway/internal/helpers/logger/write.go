package logger

// using deffrent method's for see log messages on,
// deffrent labels

func (l *Logger) Info(content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[INFO] %s", content)
}

func (l *Logger) Error(content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[ERROR] %s", content)
}

func (l *Logger) Warn(content string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Printf("[WARNING] %s", content)
}
