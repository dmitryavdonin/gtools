package logger

import (
	"fmt"
	"os"
	"strings"

	formatters "github.com/fabienm/go-logrus-formatters"
	"gitlab.com/kanya384/gotools/context"

	"github.com/sirupsen/logrus"
)

// Interface -.
type Interface interface {
	Debug(message interface{}, args ...interface{})
	DebugWithContext(ctx context.Context, message string, args ...interface{})
	Info(message string, args ...interface{})
	InfoWithContext(ctx context.Context, message string, args ...interface{})
	Warn(message string, args ...interface{})
	WarnWithContext(ctx context.Context, message string, args ...interface{})
	Error(message string, args ...interface{})
	ErrorWithContext(ctx context.Context, message string, args ...interface{})
	Fatal(message interface{}, args ...interface{})
	FatalWithContext(ctx context.Context, message string, args ...interface{})
}

// Logger -.
type Logger struct {
	*logrus.Logger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(level, serviceName string) (*Logger, error) {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warn":
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	logger := logrus.New()

	logger.Level = l
	logger.Formatter = formatters.NewGelf(serviceName)
	logger.Out = os.Stdout

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) getContextFields(ctx context.Context) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields{"traceID": ctx.TraceID()})
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args...)
}

func (l *Logger) DebugWithContext(ctx context.Context, message string, args ...interface{}) {
	l.getContextFields(ctx).Debugf(message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) InfoWithContext(ctx context.Context, message string, args ...interface{}) {
	l.getContextFields(ctx).Infof(message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(message, args...)
}

func (l *Logger) WarnWithContext(ctx context.Context, message string, args ...interface{}) {
	l.getContextFields(ctx).Warnf(message, args...)
}

// Error -.
func (l *Logger) Error(message string, args ...interface{}) {
	l.Logger.Errorf(message, args)
}

func (l *Logger) ErrorWithContext(ctx context.Context, message string, args ...interface{}) {
	l.getContextFields(ctx).Errorf(message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}

func (l *Logger) FatalWithContext(ctx context.Context, message string, args ...interface{}) {
	l.getContextFields(ctx).Fatalf(message, args...)
}

func (l *Logger) log(message string, args ...interface{}) {
	if len(args) == 0 {
		l.Logger.Info(message)
	} else {
		l.Logger.Infof(message, args...)
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(msg.Error(), args...)
	case string:
		l.log(msg, args...)
	default:
		l.log(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}
