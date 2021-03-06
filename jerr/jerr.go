package jerr

import (
	"fmt"
	"strings"
)

type JError struct {
	Messages []string
}

const (
	boldStart    = "\x1b[1m"
	formatEnd    = "\x1b[0m"
	colorDefault = "\x1b[39m"
	colorRed     = "\x1b[31m"
	colorYellow  = "\x1b[33m"
)

func (e JError) Error() string {
	return e.getText(false)
}

func (e JError) Print() {
	fmt.Println(e.getText(false))
}

func (e JError) Warn() {
	fmt.Println(e.getText(true))
}

func (e JError) getText(warn bool) string {
	returnString := ""
	for i := len(e.Messages) - 1; i >= 0; i-- {
		returnString += "\n " + boldStart + "[" + fmt.Sprintf("%d", len(e.Messages)-i) + "]" + formatEnd + " " + e.Messages[i]
	}
	if warn {
		return boldStart + colorYellow + "Warning:" + colorDefault + formatEnd + returnString
	} else {
		return boldStart + colorRed + "Error:" + colorDefault + formatEnd + returnString
	}
}

func Get(message string, err error) JError {
	var returnError JError
	switch t := err.(type) {
	case JError:
		returnError = t
		returnError.Messages = append(returnError.Messages, message)
	default:
		returnError = JError{
			Messages: []string{
				err.Error(),
				message,
			},
		}
	}
	return returnError
}

func Getf(err error, format string, a ...interface{}) JError {
	return Get(fmt.Sprintf(format, a...), err)
}

func New(message string) JError {
	return JError{
		Messages: []string{message},
	}
}

func Newf(format string, a ...interface{}) JError {
	return New(fmt.Sprintf(format, a...))
}

func Combine(errorArray ...error) error {
	var returnError JError
	for _, err := range errorArray {
		switch t := err.(type) {
		case JError:
			returnError.Messages = append(returnError.Messages, t.Messages...)
		default:
			returnError.Messages = append(returnError.Messages, t.Error())
		}
	}
	return returnError
}

func HasError(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(JError)
	if !ok {
		return e.Error() == s
	}
	for _, errMessage := range err.Messages {
		if errMessage == s {
			return true
		}
	}
	return false
}

func HasErrorPart(e error, s string) bool {
	if e == nil {
		return false
	}
	err, ok := e.(JError)
	if !ok {
		return strings.Contains(e.Error(), s)
	}
	for _, errMessage := range err.Messages {
		if strings.Contains(errMessage, s) {
			return true
		}
	}
	return false
}
