package logger

import (
	"io"
	"log"
	"os"
)

func GetLogger(verbose bool) (*log.Logger, func()) {
	logFile, err := os.OpenFile("clai.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}

	var output io.Writer

	if verbose {
		output = io.MultiWriter(os.Stdout, logFile)
	} else {
		output = logFile
	}

	logger := log.New(output, "", log.Ldate|log.Ltime|log.Lshortfile)

	return logger, func() {
		logFile.Close()
	}
}
