package main

import (
	"sync"
	"time"
)

// -------------------------- money system ----------------------- \\
var wallet int64
var bank int64

// ------------------------------ game -------------------------------\\
var save string
var hungry int
var TOWN string
var PLACE string
var NAME string
var autosave int = -1

// ------------------------------ inputs handle --------------------------- \\

var state bool
var mu sync.Mutex

func setState(newState bool) {
	mu.Lock()
	state = newState
	mu.Unlock()
}

func getState() bool {
	mu.Lock()
	defer mu.Unlock()
	return state
}

// --------------------------- time ----------------------------------\\

var ms = time.Millisecond
var sec = time.Second
var started bool = false

// ------------------------ text formating --------------------------- \\
var colorCodes = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}
