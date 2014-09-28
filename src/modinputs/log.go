package modinputs

import "os"
import "fmt"

const (
	LEVEL_DEBUG = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

type log struct {
	Level int
}

func (l log) Debug(msg string, v ...interface{}) {
	if l.Level <= LEVEL_DEBUG {
		fmt.Fprintf(os.Stderr, "DEBUG %s\n", fmt.Sprintf(msg, v...))
	}
}

func (l log) Info(msg string, v ...interface{}) {
	if l.Level <= LEVEL_INFO {
		fmt.Fprintf(os.Stderr, "INFO %s\n", fmt.Sprintf(msg, v...))
	}
}

func (l log) Warn(msg string, v ...interface{}) {
	if l.Level <= LEVEL_WARN {
		fmt.Fprintf(os.Stderr, "WARN %s\n", fmt.Sprintf(msg, v...))
	}
}

func (l log) Error(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR %s\n", fmt.Sprintf(msg, v...))
}

func (l log) Fatal(msg string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL %s\n", fmt.Sprintf(msg, v...))
}

var Log *log = &log{Level: LEVEL_INFO}
