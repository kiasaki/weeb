package weeb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// L is a shorthand type for log message fields and context
type L map[string]interface{}

const (
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
	LogLevelFatal   = "fatal"
)

var logLevelsOrder = map[string]int{
	LogLevelDebug:   1,
	LogLevelInfo:    2,
	LogLevelWarning: 3,
	LogLevelError:   4,
	LogLevelFatal:   5,
}

var currentLogLevel = LogLevelDebug

// SetGlobalLogLevel sets the global log level
func SetGlobalLogLevel(level string) {
	currentLogLevel = level
}

// Logger is a logger instance
type Logger struct {
	formatter func(L) string
	outputs   []func(string)
	context   L
}

// NewLogger creates a new Logger instance
func NewLogger() *Logger {
	return &Logger{
		formatter: defaultLogFormatter,
		outputs:   []func(string){defaultLogOutput},
		context:   L{},
	}
}

// SetFormatter sets the formatter function for this logger
func (l *Logger) SetFormatter(fn func(L) string) {
	l.formatter = fn
}

// ClearOutputs removes all the output functions from this logger
func (l *Logger) ClearOutputs(fn func(string)) {
	l.outputs = []func(string){}
}

// AddOutput adds a new output function to this logger
func (l *Logger) AddOutput(fn func(string)) {
	l.outputs = append(l.outputs, fn)
}

// SetContext merge extra context values into this logger's instance
func (l *Logger) SetContext(addedContext L) *Logger {
	for k, v := range addedContext {
		l.context[k] = v
	}
	return l
}

// WithContext create a new logger sub instance with added context
func (l *Logger) WithContext(addedContext L) *Logger {
	newLogger := NewLogger()
	newLogger.formatter = l.formatter
	newLogger.outputs = l.outputs
	for k, v := range l.context {
		newLogger.context[k] = v
	}
	for k, v := range addedContext {
		newLogger.context[k] = v
	}
	return newLogger
}

// Log sends a formatted log message to all configured outputs
func (l *Logger) Log(level, msg string, extra L) {
	// Abort if level is lower than current log level
	if logLevelsOrder[level] < logLevelsOrder[currentLogLevel] {
		return
	}

	message := L{}
	for k, v := range l.context {
		message[k] = v
	}
	for k, v := range extra {
		message[k] = v
	}
	message["msg"] = msg
	message["time"] = time.Now().UTC().Format(time.RFC3339)
	message["level"] = level
	messageString := l.formatter(message) + "\n"
	for _, output := range l.outputs {
		output(messageString)
	}
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, extra L) {
	l.Log(LogLevelDebug, msg, extra)
}

// Info logs a debug level message
func (l *Logger) Info(msg string, extra L) {
	l.Log(LogLevelInfo, msg, extra)
}

// Warning logs a debug level message
func (l *Logger) Warning(msg string, extra L) {
	l.Log(LogLevelWarning, msg, extra)
}

// Error logs a debug level message
func (l *Logger) Error(msg string, extra L) {
	l.Log(LogLevelError, msg, extra)
}

// Fatal logs a debug level message
func (l *Logger) Fatal(msg string, extra L) {
	l.Log(LogLevelFatal, msg, extra)
}

// Println logs and info level message
func (l *Logger) Println(values ...interface{}) {
	msg := fmt.Sprintf("%v", values)
	l.Info(msg[1:len(msg)-1], L{})
}

func defaultLogFormatter(message L) string {
	if bytes, err := json.Marshal(message); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

func defaultLogOutput(message string) {
	if terminal.IsTerminal(int(os.Stdout.Fd())) && message[0] == '{' && message[len(message)-2] == '}' {
		prettyLogOutput(message)
	} else {
		fmt.Print(message)
	}
}

func prettyLogOutput(message string) {
	value := L{}
	if err := json.Unmarshal([]byte(message), &value); err != nil {
		panic(err)
	}

	parsedTime, err := time.Parse(time.RFC3339, value["time"].(string))
	if err != nil {
		panic(err)
	}
	prettyTime := parsedTime.Format("2006/01/02 15:06:07")

	level := strings.ToUpper(value["level"].(string))
	fmt.Printf("%v %v %v ", prettyTime, level, value["msg"])

	delete(value, "time")
	delete(value, "level")
	delete(value, "msg")
	args, err := json.Marshal(&value)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(args))
}
