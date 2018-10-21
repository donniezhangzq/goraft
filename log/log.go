package log

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"io"
	"sync"
)

type Logger struct {
	logger *log.Logger
	entry  *log.Entry
	mu     *sync.Mutex
}

func NewLogger() *Logger {
	logger := log.New()
	entry := log.NewEntry(logger)
	return &Logger{
		logger: logger,
		entry:  entry,
	}
}

func (l *Logger) InitLogger(options *Options) error {
	file, err := l.createLogPath(options.LogPath)
	if err != nil {
		return err
	}
	l.SetOutput(file)
	l.SetFormatter(&log.TextFormatter{})
	level, err := log.ParseLevel(options.LogLevel)
	if err != nil {
		return err
	}
	l.SetLevel(level)
	return nil
}

func (l *Logger) createLogPath(logPath string) (*os.File, error) {
	dir, err := filepath.Abs(filepath.Dir(logPath))
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_APPEND, 0644)
	return file, err
}

func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Out = out
}

func (l *Logger) SetFormatter(formatter log.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Formatter = formatter
}

func (l *Logger) SetLevel(level log.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Level = level
}

func (l *Logger) SetDefaultField(role ElectionState, id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entry = l.entry.WithFields(log.Fields{
		"role": role,
		"id":   id,
	})
}

func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}


type FatalHook struct {
	f      func(entry *log.Entry) error
	logger *Logger
}

func NewFatalHook(f func(entry *log.Entry) error, logger *Logger) *FatalHook {
	return &FatalHook{
		f:      f,
		logger: logger,
	}
}

func (fh *FatalHook) AddHook(hook log.Hook) {
	fh.logger.mu.Lock()
	defer fh.logger.mu.Unlock()
	for _, level := range hook.Levels() {
		fh.logger.logger.Hooks[level] = append(fh.logger.logger.Hooks[level], hook)
	}
}

func (fh FatalHook) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
	log.AddHook()
}

func (fh FatalHook) Fire(entry *log.Entry) error {
	return fh.f(entry)
}
