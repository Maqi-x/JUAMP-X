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
