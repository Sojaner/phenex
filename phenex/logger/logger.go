package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	FilePath string
	file     *os.File
}

func Create(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return &Logger{FilePath: filePath, file: file}, err
}

func (Logger *Logger) Printf(format string, v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprintf(format, v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Println(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprintln(v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Print(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprint(v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Errorf(format string, v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ERROR: ", time.DateTime)
	message := dateTime + fmt.Sprintf(format, v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Errorln(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ERROR: ", time.DateTime)
	message := dateTime + fmt.Sprintln(v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Error(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ERROR: ", time.DateTime)
	message := dateTime + fmt.Sprint(v...)
	_, _ = Logger.file.WriteString(message)
	log.Print(message)
}

func (Logger *Logger) Fatal(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ERROR: ", time.DateTime)
	message := dateTime + fmt.Sprint(v...)
	_, _ = Logger.file.WriteString(message)
	_ = Logger.file.Close()
	log.Fatal(message)
}
