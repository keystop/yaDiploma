package logger

import (
	"log"
	"os"
	"sync"
)

var infoLog *log.Logger
var errorLog *log.Logger
var fInfo *os.File
var fError *os.File
var once sync.Once

func Close() {
	fInfo.Close()
}

func Infof(str string, v ...interface{}) {
	infoLog.Printf(str, v...)
}

func Info(v ...interface{}) {
	infoLog.Println(v...)
}

func Error(v ...interface{}) {
	infoLog.Println("ERROR")
	infoLog.Println(v...)
	errorLog.Fatal(v...)
}

func Panic(v ...interface{}) {
	infoLog.Println("PANIC")
	infoLog.Println(v...)
	errorLog.Panic(v...)
}

func createLogs() {
	var err error
	fInfo, err = os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	infoLog = log.New(fInfo, "INFO\t", log.Ldate|log.Ltime)

	fError, err = os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	errorLog = log.New(fError, "ERROR\t", log.Ldate|log.Ltime)
}

func NewLogs() {
	once.Do(createLogs)
}
