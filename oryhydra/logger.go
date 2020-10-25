package oryhydra

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	lg logger = nopLogger{}
)

func init() {
	fpath := os.Getenv("TERRAFORM_PROVIDER_ORYHYDRA_LOG_FILE_PATH")
	if fpath == "" {
		return
	}

	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("open log file: %v", err))
	}

	lg = log.New(f, "", 0)
	lg.Printf("----- %s -----", time.Now().Format(time.RFC3339))

	log.SetOutput(f)
}

type logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type nopLogger struct{}

func (n nopLogger) Print(v ...interface{}) {}

func (n nopLogger) Printf(format string, v ...interface{}) {}
