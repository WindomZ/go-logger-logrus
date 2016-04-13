package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

type Logger struct {
	logger         logrus.Logger
	_entry         *Entry
	dev            bool
	writer         *fileWriter
	writeFormatter logrus.Formatter
	callDepth      int
	callTemp       int
}

func New(c *Config) *Logger {
	s := &Logger{logger: *logrus.New()}
	s.dev = c.Dev
	s.logger.Formatter = getFormatter(c.Formatter)
	s.logger.Level = c.Level
	if c.CallDepth > 0 {
		s.callDepth = c.CallDepth
	} else {
		s.callDepth = 0
	}
	s.callTemp = 0
	if len(c.FilePath) != 0 {
		s.writer = newFileWriter(c)
		s.writeFormatter = getFormatter(c.FileFormatter)
		err := s.writer.Init()
		if err != nil {
			panic(err)
		}
	} else {
		s.writer = nil
	}
	return s
}

func (logger *Logger) IsDebug() bool {
	return logger.dev
}

func (logger *Logger) pass(level logrus.Level) bool {
	return logger.logger.Level >= level
}

func (logger *Logger) entry() *Entry {
	if logger._entry == nil {
		logger._entry = NewEntry(logger)
		if logger.callTemp != 0 {
			logger._entry.call(logger.callTemp)
			logger.callTemp = 0
		}
	}
	return logger._entry
}

func (logger *Logger) C(skip int) *Logger {
	if skip > 0 {
		logger.callTemp = skip
	}
	return logger
}

func (logger *Logger) KV(key string, value interface{}) *Logger {
	logger.entry().withField(key, value)
	return logger
}

func (logger *Logger) E(err error) *Logger {
	logger.entry().withError(err)
	return logger
}

// ====================================================

func (logger *Logger) Debug(args ...interface{}) {
	if logger.pass(DebugLevel) {
		logger.entry().Debug(args...)
	}
}

func (logger *Logger) Print(args ...interface{}) {
	logger.entry().Print(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	if logger.pass(InfoLevel) {
		logger.entry().Info(args...)
	}
}

func (logger *Logger) Warn(args ...interface{}) {
	if logger.pass(WarnLevel) {
		logger.entry().Warn(args...)
	}
}

func (logger *Logger) Error(args ...interface{}) {
	if logger.pass(ErrorLevel) {
		logger.entry().Error(args...)
	}
}

func (logger *Logger) Fatal(args ...interface{}) {
	if logger.pass(FatalLevel) {
		logger.entry().Fatal(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	if logger.pass(PanicLevel) {
		logger.entry().Panic(args...)
	}
}

// ====================================================

func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.pass(DebugLevel) {
		logger.entry().Debugf(format, args...)
	}
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.entry().Printf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.pass(InfoLevel) {
		logger.entry().Infof(format, args...)
	}
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.pass(WarnLevel) {
		logger.entry().Warnf(format, args...)
	}
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.pass(ErrorLevel) {
		logger.entry().Errorf(format, args...)
	}
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.pass(FatalLevel) {
		logger.entry().Fatalf(format, args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.pass(PanicLevel) {
		logger.entry().Panicf(format, args...)
	}
}

// ====================================================

func (logger *Logger) Debugln(args ...interface{}) {
	if logger.pass(DebugLevel) {
		logger.entry().Debugln(args...)
	}
}

func (logger *Logger) Println(args ...interface{}) {
	logger.entry().Println(args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	if logger.pass(InfoLevel) {
		logger.entry().Infoln(args...)
	}
}

func (logger *Logger) Warnln(args ...interface{}) {
	if logger.pass(WarnLevel) {
		logger.entry().Warnln(args...)
	}
}

func (logger *Logger) Errorln(args ...interface{}) {
	if logger.pass(ErrorLevel) {
		logger.entry().Errorln(args...)
	}
}

func (logger *Logger) Fatalln(args ...interface{}) {
	if logger.pass(FatalLevel) {
		logger.entry().Fatalln(args...)
	}
	os.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	if logger.pass(PanicLevel) {
		logger.entry().Panicln(args...)
	}
}
