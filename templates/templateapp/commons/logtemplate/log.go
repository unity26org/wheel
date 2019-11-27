package logtemplate

var Path = []string{"commons", "log", "log.go"}

var Content = `package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

var writer = setWriter()

var (
	Debug = log.New(writer, fmt.Sprintf("\033[94m[DEBUG]\033[39m "), log.LstdFlags)
	Info  = log.New(writer, fmt.Sprintf("\033[36m[INFO]\033[39m "), log.LstdFlags)
	Warn  = log.New(writer, fmt.Sprintf("\033[93m[WARN]\033[39m "), log.LstdFlags)
	Error = log.New(writer, fmt.Sprintf("\033[91m[ERROR]\033[39m "), log.LstdFlags)
	Fatal = log.New(writer, fmt.Sprintf("\033[90m[FATAL]\033[39m "), log.LstdFlags)
)

func setWriter() io.Writer {
	os.MkdirAll("./log", 0755)

	fileHandler, err := os.OpenFile("./log/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return io.MultiWriter(os.Stdout, fileHandler)
}`
