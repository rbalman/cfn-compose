package logger

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"

// var Blue = "\033[34m"
// var Magenta = "\033[35m"
// var Cyan = "\033[36m"

const (
	DEBUG int32 = 0
	INFO  int32 = 2
	WARN  int32 = 4
	ERROR int32 = 8
)

type Logger struct {
	Debug    *log.Logger
	Info     *log.Logger
	Warn     *log.Logger
	Error    *log.Logger
	LogLevel int32
}

var Log Logger

func Start(logLevel int32) {
	debugHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard

	switch logLevel {

	case DEBUG:
		debugHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr

	case INFO:
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr

	case WARN:
		warnHandle = os.Stdout
		errorHandle = os.Stderr

	case ERROR:
		warnHandle = os.Stdout
		errorHandle = os.Stderr

	// Defaults to INFO
	default:
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr

	}

	// Log.Debug = log.New(debugHandle, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Log.Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Log.Warn = log.New(warnHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	// Log.Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	Log.Debug = log.New(debugHandle, "[DEBUG] ", 0)
	Log.Info = log.New(infoHandle, "[INFO] ", 0)
	Log.Warn = log.New(warnHandle, "[WARNING] ", log.Lshortfile)
	Log.Error = log.New(errorHandle, "[ERROR] ", log.Lshortfile)

	atomic.StoreInt32(&Log.LogLevel, logLevel)
}

// LogLevel returns the configured logging level.
func LogLevel() int32 {
	return atomic.LoadInt32(&Log.LogLevel)
}

func (l Logger) Debugf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	l.Debug.Output(2, str)
}

func (l Logger) Debugln(params ...interface{}) {
	str := fmt.Sprintf("%s", params...)
	l.Debug.Output(2, str)
}

func (l Logger) DebugCtxf(ctx context.Context, format string, params ...interface{}) {
	ctxStr := getContextString(ctx)
	str := fmt.Sprintf(format, params...)
	l.Debug.Output(2, ctxStr+str)
}

func (l Logger) Infof(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	l.Info.Output(2, str)
}

func (l Logger) Infoln(params ...interface{}) {
	str := fmt.Sprintf("%s", params...)
	l.Info.Output(2, str)
}

func (l Logger) InfoCtxf(ctx context.Context, format string, params ...interface{}) {
	ctxStr := getContextString(ctx)
	str := fmt.Sprintf(format, params...)
	l.Info.Output(2, ctxStr+str)
}

func (l Logger) Warnf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	l.Warn.Output(2, str)
}

func (l Logger) Warnln(params ...interface{}) {
	str := fmt.Sprintf("%s", params...)
	l.Warn.Output(2, str)
}

func (l Logger) WarnCtxf(ctx context.Context, format string, params ...interface{}) {
	ctxStr := getContextString(ctx)
	str := fmt.Sprintf(format, params...)
	l.Warn.Output(2, ctxStr+str)
}

func (l Logger) Errorf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	l.Error.Output(2, str)
}

func (l Logger) Errorln(params ...interface{}) {
	str := fmt.Sprintf("%s", params...)
	l.Error.Output(2, str)
}

func (l Logger) ErrorCtxf(ctx context.Context, format string, params ...interface{}) {
	ctxStr := getContextString(ctx)
	str := fmt.Sprintf(format, params...)
	l.Error.Output(2, ctxStr+str)
}

func getContextString(ctx context.Context) (ctxStr string) {
	if job, ok := ctx.Value("job").(string); ok {
		ctxStr += fmt.Sprintf("[JOB: %s] ", job)
	}

	if stack, ok := ctx.Value("stack").(string); ok {
		ctxStr += fmt.Sprintf("[STACK: %s] ", stack)
	}

	return
}

// func (l *Logger) ColorPrint(ctx context.Context, lg string) {
// 	if color, ok := ctx.Value("logColor").(string); ok {
// 		lg = fmt.Sprintf("%s%s%s", color, lg, Reset)
// 	}
// 	fmt.Printf("%s", lg)
// }

// func (l *Logger) ColorPrintf(ctx context.Context, lg string, params ...interface{}) {
// 	str := formatString(lg, params...)
// 	l.ColorPrint(ctx, str)
// }

// func (l *Logger) Info(lg string) string {
// 	lg = fmt.Sprintf("%sℹ %s%s", Blue, lg, Reset)
// 	fmt.Printf("%s\n", lg)
// 	return lg
// }

// func (l *Logger) Infof(lg string, params ...interface{}) string {
// 	str := formatString(lg, params...)
// 	return l.Info(str)
// }

// func (l *Logger) Warn(lg string) string {
// 	lg = fmt.Sprintf("%s⚠ %s%s", Yellow, lg, Reset)
// 	fmt.Printf("%s\n", lg)
// 	return lg
// }

// func (l *Logger) Warnf(lg string, params ...interface{}) string {
// 	str := formatString(lg, params...)
// 	return l.Warn(str)
// }

// func (l *Logger) Error(lg string) string {
// 	lg = fmt.Sprintf("%s✖️ %s %s", Red, lg, Reset)
// 	fmt.Printf("%s\n", lg)
// 	return lg
// }
// func (l *Logger) Errorf(lg string, params ...interface{}) string {
// 	str := formatString(lg, params...)
// 	return l.Error(str)
// }

// func (l *Logger) Success(lg string) string {
// 	lg = fmt.Sprintf("%s✔ %s %s", Green, lg, Reset)
// 	fmt.Printf("%s\n", lg)
// 	return lg
// }

// func (l *Logger) Successf(lg string, params ...interface{}) string {
// 	str := formatString(lg, params...)
// 	return l.Success(str)
// }
// func (l *Logger) Print(lg string) string {
// 	lg = fmt.Sprintf("%s %s %s", Reset, lg, Reset)
// 	fmt.Printf("%s\n", lg)
// 	return lg
// }
// func (l *Logger) Printf(lg string, params ...interface{}) string {
// 	str := formatString(lg, params...)
// 	return l.Print(str)
// }
