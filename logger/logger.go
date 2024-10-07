package logger

import (
	"io"
	"log"
	"os"
)

func GetLogger(verbose bool) *log.Logger {
	if !verbose {
		// If not verbose, return a logger that writes to a no-op writer
		return log.New(io.Discard, "", 0) // io.Discard discards all output
	}

	// If verbose, log to stdout
	return log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
