package logrus

import "github.com/Sirupsen/logrus"

type Format uint8

const (
	FormatDefault Format = iota
	FormatJSON
	FormatLogstash
)

func getFormatter(f Format) logrus.Formatter {
	switch f {
	case FormatJSON:
		return new(logrus.JSONFormatter)
	case FormatLogstash:
	}
	return new(DefaultFormatter)
}
