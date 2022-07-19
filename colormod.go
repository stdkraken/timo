package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/color"
)

// for removing Reset in gookit/color

func RenderString(code string, str string) string {
	if len(code) == 0 || str == "" {
		return str
	}

	if !color.Enable || !color.SupportColor() {
		return color.ClearCode(str)
	}

	return color.StartSet + code + "m" + str
}

func doPrintV2(code, str string) {
	_, err := fmt.Fprint(os.Stdout, RenderString(code, str))
	saveInternalError(err)
}

func doPrintlnV2(code string, args []interface{}) {
	str := formatArgsForPrintln(args)
	_, err := fmt.Fprintln(os.Stdout, RenderString(code, str))
	saveInternalError(err)
}

func saveInternalError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func formatArgsForPrintln(args []interface{}) (message string) {
	if ln := len(args); ln == 0 {
		message = ""
	} else if ln == 1 {
		message = fmt.Sprint(args[0])
	} else {
		message = fmt.Sprintln(args...)
		message = message[:len(message)-1]
	}
	return
}

func PrintColorMod(s color.Style, a ...interface{}) {
	doPrintV2(s.String(), fmt.Sprint(a...))
}

func PrintfColorMod(s color.Style, format string, a ...interface{}) {
	doPrintV2(s.Code(), fmt.Sprintf(format, a...))
}
