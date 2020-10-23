package log

import (
	"log"
	"os"
)

var logger *log.Logger

func Init() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Fatal(err error) {
	logger.Fatal(err)
}

func Panic(err error) {
	logger.Panic(err)
}
