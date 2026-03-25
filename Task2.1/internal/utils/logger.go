package utils

import (
	"fmt"
	"log"
	"os"
)

// LogWithLabelAndTimeStamp prints a labeled message to stdout.
func LogWithLabelAndTimeStamp(label, text string) {
	l := log.New(os.Stdout, fmt.Sprintf("[%s]: ", label), log.Lmsgprefix|log.Ldate|log.Ltime)
	l.Print(text)
}
