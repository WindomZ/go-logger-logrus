package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	logrus.Entry
	logger    *Logger
	callDepth int
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{*logrus.NewEntry(&logger.logger), logger, logger.callDepth}
}

func (entry *Entry) call(skip int) *Entry {
	if skip < 0 {
		skip = 0
	}
	entry.callDepth = skip
	return entry
}

func (entry *Entry) hook() *Entry {
	if !entry.logger.dev {
		return entry
	}
	_, file, line, ok := runtime.Caller(entry.callDepth)
	if !ok {
		file = "???"
		line = 0
	}
	dir, filename := path.Split(file)
	//dir = dir[strings.Index(dir, "/src/")+5:]
	dir = dir[strings.LastIndex(dir[:strings.LastIndex(dir, "/")], "/")+1:]
	entry.Entry = *entry.WithFields(logrus.Fields{
		"dir":      dir,
		"filename": filename,
		"line":     strconv.FormatInt(int64(line), 10),
	})
	return entry
}

func (entry *Entry) withError(err error) *Entry {
	return entry.withField(logrus.ErrorKey, err)
}

func (entry *Entry) withField(key string, value interface{}) *Entry {
	return entry.withFields(logrus.Fields{key: value})
}

func (entry *Entry) withFields(fields logrus.Fields) *Entry {
	entry.Entry = *entry.WithFields(fields)
	return entry
}

func (entry *Entry) write(level logrus.Level, args ...interface{}) *Entry {
	if entry.logger.writer == nil {
		return entry
	}
	if entry.Logger.Level < level {
		return entry
	}
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = fmt.Sprint(args...)

	serialized, err := entry.logger.writeFormatter.Format(&entry.Entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	}

	err = entry.logger.writer.WriteMsg(string(serialized))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}

	return entry
}

// ====================================================

func (entry *Entry) debug(args ...interface{}) {
	entry.hook().write(DebugLevel, args...).Entry.Debug(args...)
}

func (entry *Entry) print(args ...interface{}) {
	entry.hook().Entry.Info(args...)
}

func (entry *Entry) info(args ...interface{}) {
	entry.hook().write(InfoLevel, args...).Entry.Info(args...)
}

func (entry *Entry) warn(args ...interface{}) {
	entry.hook().write(WarnLevel, args...).Entry.Warn(args...)
}

func (entry *Entry) error(args ...interface{}) {
	entry.hook().write(ErrorLevel, args...).Entry.Error(args...)
}

func (entry *Entry) fatal(args ...interface{}) {
	entry.hook().write(FatalLevel, args...).Entry.Fatal(args...)
}

func (entry *Entry) panic(args ...interface{}) {
	entry.hook().write(PanicLevel, args...).Entry.Panic(args...)
}

// ====================================================

func (entry *Entry) Debug(args ...interface{}) {
	entry.debug(args...)
}

func (entry *Entry) Print(args ...interface{}) {
	entry.print(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.info(args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.error(args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.fatal(args...)
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.panic(args...)
}

// ====================================================

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.debug(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.print(fmt.Sprintf(format, args...))
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.info(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.warn(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.error(fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.fatal(fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.panic(fmt.Sprintf(format, args...))
	}
}

// ====================================================

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.debug(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	entry.print(entry.sprintlnn(args...))
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.info(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.warn(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.error(entry.sprintlnn(args...))
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.fatal(entry.sprintlnn(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.panic(entry.sprintlnn(args...))
	}
}

func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
