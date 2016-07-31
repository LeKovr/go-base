// Package logger incapsulates logger logic
// This is a proxy for github.com/Sirupsen/logrus
// with some additions like
// functional options and logging to file
package logger

import (
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Log is a copy of logrus.Entry plus, when logging to file used, its filehandle.
type Log struct {
	*logrus.Entry
	file *os.File
}

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	Dest  string `long:"log_dest"  description:"Log destination (STDERR)"`
	Level string `long:"log_level" description:"Log level [warn|info|debug]" default:"warn"`
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

// TimeStamp adds timestamp in log output
func TimeStamp(logger *Log) error {
	return logger.setTimeStamp()
}

// Disable turns logging off
func Disable(logger *Log) error {
	return logger.setDisabled()
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

func (logger *Log) setTimeStamp() error {
	logger.Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05.000"}
	return nil
}

func (logger *Log) setDisabled() error {
	logger.Logger.Out = ioutil.Discard
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
