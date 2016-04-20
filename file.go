package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const SUFFIX string = "log"

type fileWriter struct {
	sync.Mutex

	Filename   string `json:"filename"`
	fileWriter *os.File

	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int

	Rotate bool `json:"rotate"`

	Perm os.FileMode `json:"perm"`
}

func newFileWriter(c *Config) *fileWriter {
	w := &fileWriter{
		Filename: "",
		Daily:    true,
		MaxDays:  7,
		Rotate:   true,
		Perm:     0660,
	}
	if c != nil {
		w.Filename = c.FilePath
		if c.KeepFileDays > 0 {
			w.MaxDays = c.KeepFileDays
		} else {
			w.Daily = false
		}
	}
	return w
}

func (w *fileWriter) Init() error {
	if len(w.Filename) == 0 {
		return errors.New("Must have filename")
	}
	return w.startLogger()
}

func (w *fileWriter) startLogger() error {
	file, err := w.createLogFile()
	if err != nil {
		return err
	}
	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
	w.fileWriter = file
	return w.initFd()
}

func (w *fileWriter) needRotate(day int) bool {
	return (w.Daily && day != w.dailyOpenDate)
}

func (w *fileWriter) WriteMsg(msg string) error {
	if w.Rotate {
		_, _, d := time.Now().Date()
		if w.needRotate(d) {
			w.Lock()
			if w.needRotate(d) {
				if err := w.doRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
				}
			}
			w.Unlock()
		}
	}

	w.Lock()
	_, err := w.fileWriter.Write([]byte(msg))
	w.Unlock()
	return err
}

func (w *fileWriter) createLogFile() (*os.File, error) {
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, w.Perm)
	return fd, err
}

func (w *fileWriter) initFd() error {
	fd := w.fileWriter
	_, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.dailyOpenDate = time.Now().Day()
	return nil
}

func (w *fileWriter) doRotate() error {
	_, err := os.Lstat(w.Filename)
	if err != nil {
		return err
	}

	fName := strings.TrimSuffix(w.Filename, filepath.Ext(w.Filename)) + fmt.Sprintf(".%s.%s", time.Now().Format("2006-01-02"), SUFFIX)

	w.fileWriter.Close()

	renameErr := os.Rename(w.Filename, fName)
	startLoggerErr := w.startLogger()
	go w.deleteOldLog()

	if startLoggerErr != nil {
		return fmt.Errorf("Rotate StartLogger: %s\n", startLoggerErr)
	} else if renameErr != nil {
		return fmt.Errorf("Rotate: %s\n", renameErr)
	}
	return nil
}

func (w *fileWriter) deleteOldLog() {
	dir := filepath.Dir(w.Filename)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.MaxDays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.Filename)) {
				os.Remove(path)
			}
		}
		return
	})
}

func (w *fileWriter) Destroy() {
	w.fileWriter.Close()
}

func (w *fileWriter) Flush() {
	w.fileWriter.Sync()
}
