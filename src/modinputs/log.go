package modinputs

import "os"
import "fmt"

func Debug(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "DEBUG %s\n", fmt.Sprintf(msg, v...))
}

func Info(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "INFO %s\n", fmt.Sprintf(msg, v...))
}

func Warn(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "WARN %s\n", fmt.Sprintf(msg, v...))
}

func Error(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR %s\n", fmt.Sprintf(msg, v...))
}

func Fatal(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL %s\n", fmt.Sprintf(msg, v...))
}
