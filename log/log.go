package log

import (
	"donniezhangzq/goraft/constant"
	logr "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type Logger struct {
	logger *logr.Logger
	entry  *logr.Entry
	mu     *sync.Mutex
}

func NewLogger() *Logger {
	logger := logr.New()
	entry := logr.NewEntry(logger)
	return &Logger{
		logger: logger,
		entry:  entry,
	}
}

func (l *Logger) InitLogger(logPath, logLevel string) error {
	file, err := l.createLogPath(logPath)
	if err != nil {
		return err
	}
	l.SetOutput(file)
	l.SetFormatter(&logr.TextFormatter{})
	level, err := logr.ParseLevel(logLevel)
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

func (l *Logger) SetFormatter(formatter logr.Formatter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Formatter = formatter
}

func (l *Logger) SetLevel(level logr.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Level = level
}

func (l *Logger) SetDefaultField(role constant.ElectionState, id string, leader string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entry = l.entry.WithFields(logr.Fields{
		"leader": leader,
		"role":   role,
		"id":     id,
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
	f      func(entry *logr.Entry) error
	logger *Logger
}

func NewFatalHook(f func(entry *logr.Entry) error, logger *Logger) *FatalHook {
	return &FatalHook{
		f:      f,
		logger: logger,
	}
}

func (fh *FatalHook) AddHook(hook logr.Hook) {
	fh.logger.mu.Lock()
	defer fh.logger.mu.Unlock()
	for _, level := range hook.Levels() {
		fh.logger.logger.Hooks[level] = append(fh.logger.logger.Hooks[level], hook)
	}
}

func (fh FatalHook) Levels() []logr.Level {
	return []logr.Level{logr.FatalLevel}
}

func (fh FatalHook) Fire(entry *logr.Entry) error {
	return fh.f(entry)
}
