package main

import (
	"os"

	"github.com/yowcow/go-simplelog"
)

func main() {
	logger := simplelog.New(os.Stdout, "[example] ", 2)
	logger.SetLevel(simplelog.Info)

	logger.Debug("this", "won't", "be", "logged")
	logger.Info("this", "will", "be", "logged")
	logger.Errorf("this %s be logged", "will")
}
