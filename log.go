package helper

import (
	"fmt"
	"runtime"
)

var (
	ErrorColor   = "\033[31m"
	SuccessColor = "\033[32m"
	InfoColor    = "\033[36m"
	ResetCOlor   = "\033[0m"
)

type Logger struct {
}

func init() {
	if runtime.GOOS == "windows" {
		ErrorColor = ""
		SuccessColor = ""
		InfoColor = ""
		ResetCOlor = ""
	}
}

func (l *Logger) LogInfo(args ...any) {
	fmt.Println(InfoColor, args, ResetCOlor)
}

func (l *Logger) LogError(args ...any) {
	fmt.Println(ErrorColor, args, ResetCOlor)
}
func (l *Logger) LogSuccess(args ...any) {
	fmt.Println(SuccessColor, args, ResetCOlor)
}
