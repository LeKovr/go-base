// Package logger incapsulates logger logic
// This is a proxy for github.com/Sirupsen/logrus
// with some additions like
// functional options and logging to file
package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Log is a copy of logrus.Entry plus, when logging to file used, its filehandle.
type Log struct {
	*logrus.Entry
	file *os.File
}

// Options is a program flags sample
// in form ready for use with github.com/jessevdk/go-flags
// It can be used in global flag struct:
//   type options struct {
//     logger.Options
//     ...other flags...
//   }
type Options struct {
	LogDest  string `short:"d" long:"dest" description:"Log destination (STDERR)"`
	LogLevel string `short:"l" long:"level" description:"Log level [warn|info|debug]" default:"warn"`
}

// WithFields adds a map into logger output
func (logger *Log) WithFields(fields Fields) *Log {
	l := Log{Entry: logger.Entry.WithFields(logrus.Fields(fields))}
	return &l
}

// WithField add a var/value pair into logger output
func (logger *Log) WithField(key string, value interface{}) *Log {
	l := Log{Entry: logger.Entry.WithField(key, value)}
	return &l
}

// -----------------------------------------------------------------------------
// Functional options

// Dest sets log destionation (STDERR if empty and file if given)
func Dest(s string) func(logger *Log) error {
	return func(logger *Log) error {
		return logger.setDest(s)
	}
}

// Level sets log level from a string value (debug/info/warn/error/fatal/panic)
func Level(s string) func(logger *Log) error {
	return func(logger *Log) error {
		return logger.setLevel(s)
	}
}

// -----------------------------------------------------------------------------
// Internal setters

func (logger *Log) setDest(s string) error {
	if s == "" {
		logger.Logger.Out = os.Stderr
	} else {
		w, err := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		logger.Logger.Out = w
		logger.file = w
		logger.Logger.Formatter = &logrus.TextFormatter{DisableColors: true, TimestampFormat: "2006-01-02 15:04:05.000"}
	}
	return nil
}

func (logger *Log) setLevel(s string) error {
	lev, err := logrus.ParseLevel(s)
	if err != nil {
		return err
	}
	logger.Logger.Level = lev
	return nil
}

// -----------------------------------------------------------------------------

// New creates a logger object
// Configuration should be set via functional options
func New(options ...func(logger *Log) error) (*Log, error) {
	l := Log{Entry: logrus.NewEntry(logrus.New())}
	l.Logger.Level = l.Level // panic
	for _, option := range options {
		err := option(&l)
		if err != nil {
			return nil, err
		}
	}
	return &l, nil
}

// -----------------------------------------------------------------------------

// Close closes filehandle if it was used
func (logger *Log) Close() error {
	if logger.file != nil {
		logger.Debug("Closing logfile")
		return logger.file.Close()
	}
	return nil
}

//pc, file, line, _ := runtime.Caller(2)
