package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var info interface{}

func handleKasyno() {
	if age < 16 {
		var odp bool
		Talk([][2]interface{}{
			{"Obsługa", "Co ty tu robisz gówniarzu??"},
			{"Obsługa", "To nie miejsce dla dzieci, wypier*alaj!"}, // lmao
			{"Ty", func() string {
				commands("Przepraszam, już idę", "Słucham??")
				for {
					inp := Prompt(">>> ")
					switch inp {
					case "1":
						return "Przepraszam, już idę"
					case "2":
						odp = true
						return "Słucham??"
					default:
						PrintClr("Niepoprawna opcja", "red")
						continue
					}
				}
			}},
		}, map[string]string{
			"Obsługa": "magenta",
			"Ty":      "green",
		})
		if odp {
			PrintClr("Zostałeś wyjebany z kasyna 😇", "red")
		}
		loading(1, "")
		Back()
	} else {
		Talk([][2]interface{}{
			{"Obsługa", "Witamy ponownie!"},
			{"Obsługa", "podobno dziś szczęśliwy dzień, może będzie jakaś wygrana?"},
			{"Ty", "Hm..."},
			{"Ty", "Oby!"},
		}, map[string]string{
			"Obsługa": "magenta",
			"Ty":      "green",
		})
		for {
			commands("Ruletka", "Blackjack", "Wyjście (nie zalecane, lepiej wyjebać znaczy zainwestować więcej)")
			inp := Prompt(">>> ")
			switch inp {
			case "1":
				handleRuletka()
			case "2":
				handleBj()
			case "3":
				Back()
			default:
				Println("Niepoprawna opcja")
				continue
			}
			chm("KASYNO")
		}
	}
}

func handleRuletka() {
	var c float64
	Talk([][2]interface{}{
		{"Kupier", "Witaj ponownie!"},
		{"Kupier", "To za ile dziś gramy?"},
		{"Ty", func() string {
			Println("Podaj liczbe złotych lub exit aby wyjść")
			for {
				inp := Prompt(">>> ")
				if inp == "exit" {
					Back()
				}
				if idk, err := strconv.ParseFloat(inp, 64); err == nil && idk > 0 {
					if wallet >= func() float64 {
						c = idk
						return idk
					}() {
						return inp
					} else {
						ShowError("Nie masz wystaraczającej ilości pieniędzy w portfelu")
						continue
					}
				} else {
					ShowError("Niepoprawna liczba")
					continue
				}
			}
		}},
		{"Kupier", "Świetnie! liczba, kolor...?"},
		{"Ty", func() string {
			PrintLine("Kolory:")
			Println("Czerwony")
			Println("Czarny")
			Println("Zielony")
			Sep()
			Println("Możesz też wybrać konkretną liczbe poprostu ją pisząc (można wpisać 0 jeśli chcesz grać bez stawki)")
			Println("(lub przedział liczb w takim formacie: {x}-{y})")
			rg := regexp.MustCompile(`^\d+-\d+$`)
			for {
				inp := Prompt(">>> ")
				if normStr(inp) == "czerwony" || normStr(inp) == "czarny" || normStr(inp) == "zielony" {
					info = normStr(inp)
					return "No cóż... chyba postawie na " + inp + "!"
				} else if x, err := strconv.Atoi(inp); err == nil {
					if x > 36 || x < 1 {
						Println("Podaj liczbe z zakresu od 1 do 36!")
						continue
					}
					info = x
					return "Liczba! to będzie... " + inp + "!"
				} else if rg.MatchString(inp) {
					info = []int{}
					x1, err1 := strconv.Atoi(strings.Split(inp, "-")[0])
					x2, err2 := strconv.Atoi(strings.Split(inp, "-")[1])

					if err1 != nil || err2 != nil || x1 < 1 || x1 > 36 || x2 < 1 || x2 > 36 {
						Println("Podaj liczby z zakresu od 1 do 36!")
						continue
					}

					info = []int{x1, x2}
					return Sprintf("Hmm... myśle że to będzie zakres od %s do %s!", strings.Split(inp, "-")[0], strings.Split(inp, "-")[1])
				} else {
					Println("Niepoprawna opcja")
					continue
				}
			}
		}},
		{"Kupier", "Świetnie! to zaczynamy!"},
	}, map[string]string{
		"Kupier": "bright_blue",
		"Ty":     "green",
	})
	win, multi := rulette()
	chm("KASYNO")
	if win {
		PrintC(Sprintf("%sWYGRAŁEŚ!%s", colorCodes["green"], colorCodes["reset"]))
		PrintClr("\033[1mStawka: "+strconv.FormatFloat(c, 'f', 2, 64), "green")
		PrintClr("\033[1mMnożnik: "+strconv.Itoa(multi), "green")
		PrintClr("\033[1mWygrana: "+strconv.FormatFloat(c*float64(multi), 'f', 2, 64), "green")
		walletAdd(c * float64(multi))
		dalej()
	} else {
		PrintC(Sprintf("%sPRZEGRAŁEŚ!%s", colorCodes["red"], colorCodes["reset"]))
		PrintClr("\033[1mStawka: "+strconv.FormatFloat(c, 'f', 2, 64), "red")
		//PrintClr("\033[1mMnożnik: "+strconv.Itoa(multi), "red")
		PrintClr("\033[1mStrata: "+strconv.FormatFloat(c, 'f', 2, 64), "red")
		walletD(c)
		dalej()
	}
}

func rulette() (bool, int) {
	var lastclr, lastindex int
	colors := genColors()
	chm("LOSOWANIE")
	Println("")
	PrintLine("LOSOWANIE:")
	for i := 0; i < randint(23, 60); i++ {
		index := i % len(colors)
		color := colors[index]
		Printf("\r%s%s%s",
			func() string {
				switch color {
				case 1:
					return "\033[41;97m"
				case 2:
					return "\033[30;47m"
				case 3:
					return "\033[42;97m"
				}
				return ""
			}(),
			func() string {
				return strings.Repeat(" ", func() int {
					w, _ := getTerminalSize()
					return w
				}())
			}(),
			colorCodes["reset"])
		lastclr, lastindex = color, index
		time.Sleep(100 * ms)
	}
	PrintLine("Wyniki: ")
	Printfln("Kolor: %s", func() string {
		switch lastclr {
		case 1:
			return "Czerwony"
		case 2:
			return "Czarny"
		case 3:
			return "Zielony"
		}
		return ""
	}())
	Printfln("Numer: %d", lastindex)
	loading(2, "Ładowanie")
	var win bool
	var multi int
	switch info.(type) {
	case string:
		if info == "zielony" && lastclr == 3 {
			win = true
			multi = 35
		} else if info == "czerwony" && lastclr == 1 {
			win = true
			multi = 2
		} else if info == "czarny" && lastclr == 2 {
			win = true
			multi = 2
		}
	case []int:
		x := info.([]int)
		if lastindex >= x[0] && lastindex <= x[1] {
			win = true

			if x[1]-x[0]+1 > 0 && x[1]-x[0]+1 <= 36 {
				multi = 37/x[1] - x[0] + 1
			} else {
				multi = 0
			}
		}
	case int:
		if lastindex == info.(int) {
			win = true
		}
	}
	Sep()
	time.Sleep(1 * sec)
	chm("KASYNO")
	return win, multi
}

func genColors() []int {
	colors := []int{3}

	redNumbers := map[int]bool{
		1: true, 3: true, 5: true, 7: true, 9: true, 12: true, 14: true, 16: true,
		18: true, 19: true, 21: true, 23: true, 25: true, 27: true, 30: true,
		32: true, 34: true, 36: true,
	}

	for i := 1; i <= 36; i++ {
		//
		if redNumbers[i] {
			colors = append(colors, 1) // Czerwony
		} else {
			colors = append(colors, 2) // Czarny
		}
	}

	return colors
}

// ----------------------------------------------------- BLACK JACK ----------------------------------------------- \\
type Card struct {
	Suit  string
	Value string
}

// Predefined deck values and suits
var suits = []string{"Hearts", "Diamonds", "Clubs", "Spades"}
var values = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King", "Ace"}
var deck []Card

// InitializeDeck creates and returns a full deck of cards
func InitializeDeck() []Card {
	deck := []Card{}
	for _, suit := range suits {
		for _, value := range values {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

func ShuffleDeck(deck []Card) []Card {
	shuffledDeck := make([]Card, len(deck))
	perm := rand.Perm(len(deck))
	for i, v := range perm {
		shuffledDeck[v] = deck[i]
	}
	return shuffledDeck
}

func GetCardValue(card Card) int {
	switch card.Value {
	case "Ace":
		return 11
	case "King", "Queen", "Jack":
		return 10
	default:
		var value int
		fmt.Sscanf(card.Value, "%d", &value)
		return value
	}
}

// CalcHand int calculates the total value of a hand
func CalcHand(hand []Card) int {
	var aces, total int

	for _, card := range hand {
		value := GetCardValue(card)
		if card.Value == "Ace" {
			aces++
		}
		total += value
	}

	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}

	return total
}

func handleBj() {
	var c float64
	Talk([][2]interface{}{
		{"Kupier", "Witaj ponownie!"},
		{"Kupier", "Ile stawiasz?"},
		{"Ty", func() string {
			Println("Wpisz liczbe zł, exit aby uciec jak nie masz kasy xdd")
			for {
				inp := Prompt(">>> ")
				if inp == "exit" {
					Back()
				}
				if idk, err := strconv.ParseFloat(inp, 64); err == nil && idk > 0 {
					if wallet >= func() float64 {
						c = idk
						return idk
					}() {
						return inp
					} else {
						Println("Nie masz wystaraczającej ilości pieniędzy w portfelu")
						continue
					}
				} else {
					Println("Niepoprawna liczba")
					continue
				}
			}
		}},
		{"Kupier", "No to zaczynamy"},
	}, map[string]string{
		"Kupier": "bright_blue",
		"Ty":     "green",
	})
	chm("BLACKJACK")
	win := BjGame()
	chm("KASYNO")
	if win == 1 {
		Print(colorCodes["green"])
		PrintLine("Nagroda")
		loading(1, "")
		Printfln("%.2f zł", c*2)
		walletAdd(float64(c * 2))
	} else if win == 0 {
		Print(colorCodes["red"])
		PrintLine("Strata")
		loading(1, "")
		Printfln("%.2f zł", c)
		walletD(float64(c))
	} else {
		Print(colorCodes["blue"])
		Println("Remis, nic nie otrzymujesz")
	}
	Print(colorCodes["reset"])
	dalej()
}

func BjGame() int {
	ctostr := func(hand []Card) string {
		var a string
		for i, card := range hand {
			if i == 0 {
				continue
			}
			a += ", " + card.Value
		}
		return hand[0].Value + a
	}
	inf := func(clr string, playerHand []Card, playerValue int, dealerHand []Card, dealerValue int) {
		Printfln("%s\033[1m", colorCodes[clr])
		PrintLine("Wyniki:")
		Printfln("\033[1mTwoje karty:\033[22m %s, \033[1mWartość:\033[22m %d", ctostr(playerHand), playerValue)
		Printfln("\033[1mKarty krupiera:\033[22m %s, \033[1mWartość: %d\033[0m", ctostr(dealerHand), dealerValue)
		dalej()
	}
	deck := ShuffleDeck(InitializeDeck())

	var playerHand, dealerHand []Card
	playerHand = append(playerHand, deck[0], deck[1])
	dealerHand = append(dealerHand, deck[2], deck[3])
	deck = deck[4:]

	for {
		playerPassed := false
		dealerPassed := false

		// Tura gracza
		playerValue := CalcHand(playerHand)
		Println("")
		PrintColor("<bold>Twoje karty:</bold>", ctostr(playerHand))
		PrintColor("<bold>Wartość twojej ręki:</bold>", playerValue)
		PrintColor("<bold>Karta krupiera:</bold>", ctostr(dealerHand))
		PrintColor("<bold>Wartość twojej ręki:</bold>", CalcHand(dealerHand))

		if playerValue > 21 {
			Println("Przegrałeś - przekroczyłeś 21!")
			inf("red", playerHand, playerValue, dealerHand, CalcHand(dealerHand))
			return 0
		}

		Println("Co chcesz zrobić? (h)it / (p)ass")
		move := Prompt(">>> ")

		if move == "h" || move == "hit" {
			playerHand = append(playerHand, deck[0])
			deck = deck[1:]
		} else if move == "p" || move == "pass" {
			playerPassed = true
		} else {
			Println("Niepoprawna opcja.")
			continue
		}

		Println("")
		Println("Karty krupiera:", ctostr(dealerHand))
		Println("Wartość ręki krupiera:", CalcHand(dealerHand))

		if dealerShouldDraw(dealerHand) {
			loading(1, "Krupier dobiera kartę")
			dealerHand = append(dealerHand, deck[0])
			deck = deck[1:]
		} else {
			Println("Krupier postanowił zostać (pass).")
			dealerPassed = true
		}

		if CalcHand(dealerHand) > 21 {
			Println("Wygrałeś - dealer przekroczył 21!")
			inf("red", playerHand, playerValue, dealerHand, CalcHand(dealerHand))
			return 0
		}

		if playerPassed && dealerPassed {
			Println("Obaj gracze postanowili zpassować. Gra skończona!")
			break
		}
	}

	playerValue := CalcHand(playerHand)
	dealerValue := CalcHand(dealerHand)

	var clr string
	var win int

	if playerValue == 21 && dealerValue != 21 {
		clr = "green"
		PrintClr("Wygrałeś!", clr)
		win = 1
	} else if dealerValue == 21 && playerValue != 21 {
		clr = "red"
		PrintClr("Przegrałeś!", clr)
		win = 0
	} else if dealerValue > 21 || playerValue > dealerValue {
		clr = "green"
		PrintClr("Wygrałeś!", clr)
		win = 1
	} else if playerValue < dealerValue {
		clr = "red"
		win = 0
		PrintClr("Przegrałeś!", clr)
	} else {
		win = -1
		clr = "blue"
		PrintClr("Remis!", clr)
	}
	inf(clr, playerHand, playerValue, dealerHand, CalcHand(dealerHand))
	return win
}

func dealerShouldDraw(dealerHand []Card) bool {
	dealerValue := CalcHand(dealerHand)

	if dealerValue >= 17 {
		return false
	}

	var riskFactor float64
	switch {
	case dealerValue <= 11:
		riskFactor = 1
	case dealerValue == 12:
		riskFactor = 0.85
	case dealerValue == 13:
		riskFactor = 0.70
	case dealerValue == 14:
		riskFactor = 0.55
	case dealerValue == 15:
		riskFactor = 0.40
	case dealerValue == 16:
		riskFactor = 0.25
	}

	riskFactor += rand.Float64() * 0.1
	return rand.Float64() < riskFactor
}
