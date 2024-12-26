package main

import (
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
			{"Obsługa", "To nie miejsce dla dzieci, wypier*alaj!"},
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
		loading(1, "")
		if odp {
			PrintClr("Zostałeś wyjebany z kasyna 😇", "red")
		}
		time.Sleep(1 * sec)
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
				//handleBj()
			case "3":
				Back()
			default:
				Println("Niepoprawna opcja")
				continue
			}
			chm("KASYNO")
		}
	}

	Back()
}

func handleRuletka() {
	var c float64
	Talk([][2]interface{}{
		{"Kupier", "Witaj ponownie!"},
		{"Kupier", "To za ile dziś gramy?"},
		{"Ty", func() string {
			Println("Wpisz liczbe zł, exit aby uciec jak nie masz kasy xdd")
			for {
				inp := Prompt(">>> ")
				if inp == "exit" {
					Back()
				}
				if _, err := strconv.ParseFloat(inp, 64); err == nil {
					if wallet >= func() float64 {
						idk, _ := strconv.ParseFloat(inp, 64)
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
