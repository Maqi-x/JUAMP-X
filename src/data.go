package main

import (
	"math/rand"
	"sync"
	"time"
)

// -------------------------- money system ----------------------- \\
var wallet float64
var bank float64

// ------------------------------ game -------------------------------\\
/* SAVE SYSTEM */
var save string
var hungry int = 100
var TOWN string
var PLACE string
var NAME string
var autosave int = -1

/* TUTORIAL */
var first bool = true
var history []string // places history
var tutStep []int = []int{1, 0}
var age int = 18
var tmp interface{}

var assortment map[string]int // ropucha
var settings = make(map[string]interface{})
var defaultSettings = map[string]interface{}{
	"autosavetime": 15,
	"tutskip":      false,
	"colors":       true,
	"devoptions":   false,
	"theme":        "bright_cyan",
}

var devmode bool
var useSettings func()

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
	"black":          "30m",
	"red":            "91m",
	"orange":         "38;5;208m",
	"green":          "32m",
	"yellow":         "33m",
	"blue":           "34m",
	"magenta":        "35m",
	"cyan":           "36m",
	"white":          "37m",
	"bright_red":     "91m",
	"bright_green":   "92m",
	"bright_yellow":  "93m",
	"bright_blue":    "94m",
	"bright_magenta": "95m",
	"bright_cyan":    "96m",
	"x":              "96m",
	"brown":          "38;5;94m",  // brązowy
	"gray":           "90m",       // szary
	"light_gray":     "37m",       // jasnoszary
	"dark_gray":      "38;5;238m", // ciemny szary
	"light_blue":     "94m",       // jasny niebieski
	"pink":           "38;5;211m", // różowy
	"lime":           "38;5;190m", // limonkowy
	"bright_orange":  "38;5;202m", // jasny pomarańczowy
	"turquoise":      "38;5;80m",  // turkusowy
	"bold":           "1m",        // pogrubienie
	"dim":            "2m",        // przyciemnienie
	"italic":         "3m",        // kursywa (nie zawsze wspierane)
	"underline":      "4m",        // podkreślenie
	"blink":          "5m",        // miganie
	"reverse":        "7m",        // odwrócone kolory
	"hidden":         "8m",        // ukryty tekst
	"reset":          "0m",        // reset
}

func colorsCodes() {
	for k, v := range colorCodes {
		colorCodes[k] = "\033[" + v
	}
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

var talks = [][][2]interface{}{
	{
		{"Mama", "Hej, co tam?"},
		{"Ty", ""},
	},
}

func randTalk() interface{} {
	return talks[rand.Intn(len(talks))]
}

var defaultSave = SaveData{
	Player: Player{
		Hungry: 100,
		Place:  "DOM",
		Town:   "_",
		Name:   "_",
		Age:    5,
	},
	Game: Game{
		Wallet:       0,
		Bank:         1000,
		Autosave:     -1,
		First:        false,
		TutorialStep: []int{1, 0},
	},
	Assortment: Assortment{
		Buns:       11,
		Pizzas:     7,
		Newspapers: 8,
		Shoes:      5,
	},
}

var tmpses bool
