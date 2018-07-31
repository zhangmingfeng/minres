package log

import (
	"log"
	"strings"
	"os"
	"time"
	"path"
)

var suffix string

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	FilePath string
	Level    int
	handle   *log.Logger
	file     *os.File
}

func NewLogger(filePath, level string) (*Logger, error) {
	if len(filePath) == 0 {
		filePath = "./log.log"
	}
	if len(level) == 0 {
		level = "DEBUG"
	}
	var logLevel int
	switch strings.ToLower(level) {
	case "debug":
		logLevel = LevelDebug
		break
	case "info":
		logLevel = LevelInfo
		break
	case "warn":
		logLevel = LevelWarn
		break
	case "error":
		logLevel = LevelError
		break
	default:
		logLevel = LevelDebug
	}
	suffix = getSuffix()
	newFilePath := changeName(filePath, suffix)
	logFile, err := os.OpenFile(newFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return nil, err
	}
	handle := log.New(logFile, "", log.Lshortfile|log.LstdFlags)
	return &Logger{
		FilePath: filePath,
		Level:    logLevel,
		handle:   handle,
		file:     logFile,
	}, nil
}

func getSuffix() string {
	return time.Now().Format("2006-01-02")
}

func changeName(filePath, suffix string) string {
	ext := path.Ext(filePath)
	filePathNoExt := strings.TrimSuffix(filePath, ext)
	filePathNoExt = filePathNoExt + "-" + suffix
	return filePathNoExt + ext
}

func (this *Logger) Debug(v ...interface{}) {
	if this.Level > LevelDebug {
		return
	}
	this.check()
	this.handle.Println("[DEBUG]", v)
}

func (this *Logger) Info(v ...interface{}) {
	if this.Level > LevelInfo {
		return
	}
	this.check()
	this.handle.Println("[INFO]", v)
}

func (this *Logger) Warn(v ...interface{}) {
	if this.Level > LevelWarn {
		return
	}
	this.check()
	this.handle.Println("[WARN]", v)
}

func (this *Logger) Error(v ...interface{}) {
	if this.Level > LevelError {
		return
	}
	this.check()
	this.handle.Println("[ERROR]", v)
}

func (this *Logger) check() error {
	nowSuffix := getSuffix()
	if nowSuffix == suffix {
		return nil
	}
	suffix = nowSuffix
	filePath := changeName(this.FilePath, suffix)
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	this.file.Close()
	this.file = logFile
	handle := log.New(this.file, "", log.Lshortfile|log.LstdFlags)
	this.handle = handle
	return nil
}
