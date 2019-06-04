package tools

import (
	"fmt"
	"log"
)

const (
	lblDebug = "DEBUG "
	lblInfo  = "INFO "
	lblWarn  = "WARN "
	lblError = "ERROR "
	lblFatal = "FATAL "
)

func printf(level string, format string, v ...interface{}) {
	log.Printf(level+format, v...)
}
func printfln(level string, format string, v ...interface{}) {
	log.Printf(level+format+"\n", v...)
}
func println(level string, v ...interface{}) {
	log.Println(level, fmt.Sprintln(v...))
}

// Debugf - write into log with debug level
func Debugf(format string, v ...interface{}) {
	printf(lblDebug, format, v...)
}

// Debugfln - write into log with debug level
func Debugfln(format string, v ...interface{}) {
	printfln(lblDebug, format, v...)
}

// Debugln - write into log with debug level
func Debugln(v ...interface{}) {
	println(lblDebug, v...)
}

// Infof - write into log with info level
func Infof(format string, v ...interface{}) {
	printf(lblInfo, format, v...)
}

// Infofln - write into log with info level
func Infofln(format string, v ...interface{}) {
	printfln(lblInfo, format, v...)
}

// Infoln - write into log with info level
func Infoln(v ...interface{}) {
	println(lblInfo, v...)
}

//Warnf - write into log with info level
func Warnf(format string, v ...interface{}) {
	printf(lblWarn, format, v...)
}

// Warnfln - write into log with info level
func Warnfln(format string, v ...interface{}) {
	printfln(lblWarn, format, v...)
}

// Warnln - write into log with info level
func Warnln(v ...interface{}) {
	println(lblWarn, v...)
}

// Errorf - write into log with error level
func Errorf(format string, v ...interface{}) {
	printf(lblError, format, v...)
}

// Errorfln - write into log with error level
func Errorfln(format string, v ...interface{}) {
	printfln(lblError, format, v...)
}

// Errorln - write into log with error level
func Errorln(v ...interface{}) {
	println(lblError, v...)
}

// Fatalf - write into log with error level
func Fatalf(format string, v ...interface{}) {
	log.Fatalf(lblFatal+format, v...)
}

// Fatalfln - write into log with error level
func Fatalfln(format string, v ...interface{}) {
	log.Fatalf(lblFatal+format+"\n", v...)
}

// Fatalln - write into log with error level
func Fatalln(v ...interface{}) {
	log.Fatalln(lblFatal, fmt.Sprintln(v...))
}
