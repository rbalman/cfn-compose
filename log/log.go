package log

import (
	"fmt"
	"context"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
type Logger struct{}


func (l *Logger) ColorPrint(ctx context.Context, lg string) {
	if color, ok := ctx.Value("logColor").(string); ok {
		lg = fmt.Sprintf("%s%s%s", color, lg, Reset)
	}
	fmt.Printf("%s", lg)
}

func (l *Logger) ColorPrintf(ctx context.Context, lg string, params ...interface{}) {
	str := formatString(lg, params...)
	l.ColorPrint(ctx, str)
}

func (l *Logger) Info(lg string) string {
	lg = fmt.Sprintf("%sℹ %s%s", Blue, lg, Reset)
	fmt.Printf("%s\n", lg)
	return lg
}

func (l *Logger) Infof(lg string, params ...interface{}) string {
	str := formatString(lg, params...)
	return l.Info(str)
}

func (l *Logger) Warn(lg string) string {
	lg = fmt.Sprintf("%s⚠ %s%s", Yellow, lg, Reset)
	fmt.Printf("%s\n", lg)
	return lg
}

func (l *Logger) Warnf(lg string, params ...interface{}) string {
	str := formatString(lg, params...)
	return l.Warn(str)
}

func (l *Logger) Error(lg string) string {
	lg = fmt.Sprintf("%s✖️ %s %s", Red, lg, Reset)
	fmt.Printf("%s\n", lg)
	return lg
}
func (l *Logger) Errorf(lg string, params ...interface{}) string {
	str := formatString(lg, params...)
	return l.Error(str)
}

func (l *Logger) Success(lg string) string {
	lg = fmt.Sprintf("%s✔ %s %s", Green, lg, Reset)
	fmt.Printf("%s\n", lg)
	return lg
}

func (l *Logger) Successf(lg string, params ...interface{}) string {
	str := formatString(lg, params...)
	return l.Success(str)
}
func (l *Logger) Print(lg string) string {
	lg = fmt.Sprintf("%s %s %s", Reset, lg, Reset)
	fmt.Printf("%s\n", lg)
	return lg
}
func (l *Logger) Printf(lg string, params ...interface{}) string {
	str := formatString(lg, params...)
	return l.Print(str)
}

func formatString(lg string, params ...interface{}) string {
	str := fmt.Sprintf(lg, params...)
	return str
}
