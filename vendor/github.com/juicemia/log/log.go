package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger wraps *logrus.Entry.
type Logger struct {
	*logrus.Entry
}

// New returns a *Logger set up with the
// package name as a field.
func New(pkg string) *Logger {
	ent := logrus.WithFields(logrus.Fields{
		"package": pkg,
	})

	return &Logger{
		Entry: ent,
	}
}

// CloneWith returns a copy of this Logger with the passed in
// fields added to it.
func (l *Logger) CloneWith(flds map[string]interface{}) *Logger {
	ent := l.WithFields(logrus.Fields(flds))

	return &Logger{
		Entry: ent,
	}
}

// SetLevelFromEnv sets the log level by reading the
// environment variable passed in. It's case insensitive
// and maps to the following:
//
// logrus.DebugLevel: debug, trace
// logrus.InfoLevel: info
// logrus.WarnLevel: warn, warning
// logrus.ErrorLevel: error
// logrus.FatalLevel: fatal
//
// By default, if no match is found, the level is set to
// logrus.InfoLevel.
func SetLevelFromEnv(env string) {
	switch strings.ToLower(os.Getenv(env)) {
	case "debug", "trace":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}
