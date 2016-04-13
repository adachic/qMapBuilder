package main

import "fmt"

type dbgLevel int

const (
	PRODUCTION = 1 + iota
	DEBUG
	VERBOSE
)

var level dbgLevel = VERBOSE

func Dlog(format string, a ...interface{}) {
	if level < PRODUCTION {
		return;
	}
	fmt.Printf(format, a...);
}

func Dlogln(a ...interface{}) {
	if level < PRODUCTION {
		return;
	}
	fmt.Println(a...)
}

func DDlog(format string, a ...interface{}) {
	if level < DEBUG {
		return;
	}
	fmt.Printf(format, a...);
}

func DDlogln(a ...interface{}) {
	if level < DEBUG {
		return;
	}
	fmt.Println(a...)
}

func DDDlog(format string, a ...interface{}) {
	if level < VERBOSE {
		return;
	}
	fmt.Printf(format, a...);
}

func DDDlogln(a ...interface{}) {
	if level < VERBOSE {
		return;
	}
	fmt.Println(a...)
}
