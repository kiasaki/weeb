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
	l.Log("debug", msg, extra)
}

// Info logs a debug level message
func (l *Logger) Info(msg string, extra L) {
	l.Log("info", msg, extra)
}

// Warning logs a debug level message
func (l *Logger) Warning(msg string, extra L) {
	l.Log("warning", msg, extra)
}

// Error logs a debug level message
func (l *Logger) Error(msg string, extra L) {
	l.Log("error", msg, extra)
}

// Fatal logs a debug level message
func (l *Logger) Fatal(msg string, extra L) {
	l.Log("fatal", msg, extra)
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

	level := strings.ToUpper(value["level"].(string))
	fmt.Printf("%v %v %v ", value["time"], level, value["msg"])

	delete(value, "time")
	delete(value, "level")
	delete(value, "msg")
	args, err := json.Marshal(&value)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(args))
}
