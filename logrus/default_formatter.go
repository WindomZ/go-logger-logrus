package logrus

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"sort"
	"strings"
	"time"
)

const DefaultTimestampFormat = "2006/01/02 15:04:05"

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

var (
	baseTimestamp time.Time
	isTerminal    bool
)

func init() {
	baseTimestamp = time.Now()
	isTerminal = logrus.IsTerminal()
}

func miniTMS() int64 {
	ts := int64(time.Since(baseTimestamp) / time.Microsecond)
	baseTimestamp = time.Now()
	return ts
}

type DefaultFormatter struct {
	ForceColors      bool
	DisableColors    bool
	DisableTimestamp bool
	FullTimestamp    bool
	TimestampFormat  string
	DisableSorting   bool
}

func prefixFieldClashes(data logrus.Fields) {
	_, ok := data["time"]
	if ok {
		data["fields.time"] = data["time"]
	}

	_, ok = data["msg"]
	if ok {
		data["fields.msg"] = data["msg"]
	}

	_, ok = data["level"]
	if ok {
		data["fields.level"] = data["level"]
	}
}

func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	b := &bytes.Buffer{}

	prefixFieldClashes(entry.Data)

	isColorTerminal := isTerminal && (runtime.GOOS != "windows")
	isColored := (f.ForceColors || isColorTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		if !f.DisableTimestamp {
			f.appendKeyValue(b, "time", entry.Time.Format(timestampFormat))
		}
		f.appendKeyValue(b, "level", entry.Level.String())
		if entry.Message != "" {
			f.appendKeyValue(b, "msg", entry.Message)
		}
		for _, key := range keys {
			f.appendKeyValue(b, key, entry.Data[key])
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *DefaultFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())[0:1]
	if _, ok := entry.Data["dir"]; ok {
		fileText := fmt.Sprintf("\x1b[%dm%s%s", blue, entry.Data["dir"], entry.Data["filename"])
		lineText := fmt.Sprintf("\x1b[%dm%s", blue, entry.Data["line"])
		fmt.Fprintf(b, "%s [\x1b[%dm%s\x1b[0m][%09d][%s\x1b[0m:%s\x1b[0m] %s",
			entry.Time.Format(timestampFormat), levelColor, levelText, miniTMS(), fileText, lineText, entry.Message)
	} else {
		fmt.Fprintf(b, "%s [\x1b[%dm%s\x1b[0m][%09d] %s",
			entry.Time.Format(timestampFormat), levelColor, levelText, miniTMS(), entry.Message)
	}

	for _, k := range keys {
		if strings.EqualFold(k, "dir") || strings.EqualFold(k, "filename") || strings.EqualFold(k, "line") {
			continue
		}
		v := entry.Data[k]
		fmt.Fprintf(b, "\n \x1b[%dm%s\x1b[0m=%#v", levelColor, k, v)
	}
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func (f *DefaultFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {

	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}
