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

func (logger *Logger) Printf(format string, v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprintf(format, v...)
	_, _ = logger.file.WriteString(message)
	log.Print(message)
}

func (logger *Logger) Println(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprintln(v...)
	_, _ = logger.file.WriteString(message)
	log.Print(message)
}

func (logger *Logger) Print(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprint(v...)
	_, _ = logger.file.WriteString(message)
	log.Print(message)
}

func (logger *Logger) Fatal(v ...interface{}) {
	dateTime := fmt.Sprintf("[%v] ", time.DateTime)
	message := dateTime + fmt.Sprint(v...)
	_, _ = logger.file.WriteString(message)
	_ = logger.file.Close()
	log.Fatal(message)
}
