package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

type Logger struct {
	writer io.Writer
	errCh  chan introspectedError
}

func NewLogger(w io.Writer, buf int) (logger *Logger, closeFunc func()) {
	l := &Logger{
		writer: w,
		errCh:  make(chan introspectedError, buf),
	}
	l.handle()
	return l, l.Close()
}

func (l *Logger) Close() (closeFunc func()) {
	return func() { close(l.errCh) }
}

func (l *Logger) SetOutput(w io.Writer) {
	l.writer = w
}

func (l *Logger) Log(err error) {
	file, function, line := l.trace()
	l.log(err, file, function, line)
}

func (l *Logger) LogPropagate(err error) error {
	file, function, line := l.trace()
	l.log(err, file, function, line)
	return err
}

func (l *Logger) Info(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- infoLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) InfoPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- infoLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) Debug(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- debugLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) DebugPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- debugLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) Warning(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- warningLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) WarningPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- warningLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) Error(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- errorLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) ErrorPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- errorLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) Critical(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- criticalLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) CriticalPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- criticalLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) Emergency(strOrErr any) {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- emergencyLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}
}

func (l *Logger) EmergencyPropagate(strOrErr any) error {
	file, function, line := l.trace()

	err := l.error(strOrErr)

	l.errCh <- emergencyLevelError{
		Dt: time.Now(),
		Mg: err.Error(),
		Fl: file,
		Fn: function,
		Ln: line,
	}

	return err
}

func (l *Logger) handle() {
	go func() {
		for err := range l.errCh {
			j, e := json.MarshalIndent(err, "", "  ")
			if e != nil {
				_, fmterr := fmt.Fprintln(l.writer, e)
				if fmterr != nil {
					log.Println(err)
					log.Fatalln(fmterr)
				}
			} else {
				_, fmterr := fmt.Fprintln(l.writer, string(j))
				if fmterr != nil {
					log.Println(err)
					log.Fatalln(fmterr)
				}
			}
		}
	}()
}

func (l *Logger) log(e error, file string, function string, line int) {
	err, isLoggableErr := e.(LoggableError)
	if !isLoggableErr {
		l.errCh <- errorLevelError{
			Dt: time.Now(),
			Mg: e.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	}

	switch err.Level() {
	case InfoLevel:
		l.errCh <- infoLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case DebugLevel:
		l.errCh <- infoLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case WarningLevel:
		l.errCh <- warningLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case ErrorLevel:
		l.errCh <- errorLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case CriticalLevel:
		l.errCh <- criticalLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	case EmergencyLevel:
		l.errCh <- emergencyLevelError{
			Dt: time.Now(),
			Mg: err.Error(),
			Fl: file,
			Fn: function,
			Ln: line,
		}
		return
	}

	panic("logger.log(): undefined error level received")
}

func (l *Logger) trace() (file string, function string, line int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return frame.File, frame.Func.Name(), frame.Line
}

func (l *Logger) error(strOrErr any) error {
	err, isErr := strOrErr.(error)
	if isErr {
		return err
	} else {
		str, isStr := strOrErr.(string)
		if isStr {
			return errors.New(str)
		}
	}
	panic("logger.error(): logging data is not a string or error type")
}

func ToReadableLevel(err introspectedError) string {
	switch err.Level() {
	case InfoLevel:
		return InfoLevelReadable
	case DebugLevel:
		return DebugLevelReadable
	case WarningLevel:
		return WarningLevelReadable
	case ErrorLevel:
		return ErrorLevelReadable
	case CriticalLevel:
		return CriticalLevelReadable
	case EmergencyLevel:
		return EmergencyLevelReadable
	}
	panic("logger.ToReadableLevel(): received undefined error level")
}

func ToLevel(readableLevel string) int {
	switch readableLevel {
	case InfoLevelReadable:
		return InfoLevel
	case DebugLevelReadable:
		return DebugLevel
	case WarningLevelReadable:
		return WarningLevel
	case ErrorLevelReadable:
		return ErrorLevel
	case CriticalLevelReadable:
		return CriticalLevel
	case EmergencyLevelReadable:
		return EmergencyLevel
	}
	panic("logger.ToLevel(): received undefined readable level")
}