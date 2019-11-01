package kubetest

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"time"
)

var (
	success = color.GreenString("✓")
	failure = color.RedString("✗")
	writer  = os.Stdout
)

// LogMessage logs a message
func LogMessage(message string, args ...interface{}) {
	fmt.Fprintln(writer, fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(message, args...)))
}

// LogSuccess logs a success message
func LogSuccess(message string, args ...interface{}) {
	fmt.Fprintln(writer, color.GreenString(fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(message, args...))))
}

// LogFailure logs a failure message
func LogFailure(message string, args ...interface{}) {
	fmt.Fprintln(writer, color.RedString(fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), fmt.Sprintf(message, args...))))
}
