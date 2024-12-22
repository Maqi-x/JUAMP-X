package main

import (
	"sync"
	"time"
)

// -------------------------- money system ----------------------- \\
var wallet float64
var bank float64

// ------------------------------ game -------------------------------\\
var save string
var hungry int = 100
var TOWN string
var PLACE string
var NAME string
var autosave int = -1
var first bool = true
var history []string
var tutStep []int = []int{1, 0}

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
	"black":          "\033[30m",
	"red":            "\033[31m",
	"green":          "\033[32m",
	"yellow":         "\033[33m",
	"blue":           "\033[34m",
	"magenta":        "\033[35m",
	"cyan":           "\033[36m",
	"white":          "\033[37m",
	"bright_red":     "\033[91m",
	"bright_green":   "\033[92m",
	"bright_yellow":  "\033[93m",
	"bright_blue":    "\033[94m",
	"bright_magenta": "\033[95m",
	"bright_cyan":    "\033[96m",
	"orange":         "\033[38;5;215m",
	"reset":          "\033[0m",
}

// -------------------------------- Ropucha ----------------------------- \\
var products map[string]struct {
	Price       float64
	Description string
}

var productS = map[string]struct {
	Price       float64
	Description string
}{
	"Bułka": {
		Price:       2.99,
		Description: "Odnawia 20 punktów głodu, idealna na szybką przekąskę",
	},
	"Pizza": {
		Price:       5.49,
		Description: "Odnawia 40 punktów głodu, idealne na długie trasy",
	},
	"Gazeta": {
		Price:       2.00,
		Description: "Przeczytaj najnowsze informacje z %s",
	},
	"Buty": {
		Price:       19.99,
		Description: "idk", // TODO: opis lmao
	},
}

// --------------------------------- Papers ---------------------------------------- \\
var papers = []string{
	"...",
	"w przyszłości dodam to gówno, nie mam czasu tego pisać",
}
