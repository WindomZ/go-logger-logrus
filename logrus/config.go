package logrus

import "github.com/Sirupsen/logrus"

type Config struct {
	Dev           bool
	Formatter     Format
	Level         logrus.Level
	CallDepth     int
	FilePath      string
	FileFormatter Format
	KeepFileDays  int64
}

func NewDefaultConfig() *Config {
	return &Config{Dev: true, Formatter: FormatDefault, Level: DebugLevel, CallDepth: 5,
		FilePath: "", FileFormatter: FormatJSON, KeepFileDays: 7}
}
