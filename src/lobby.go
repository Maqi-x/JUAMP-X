package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func formatStep(step []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(step), " ", ",", -1), "[]")
}

func handleLobby() {
	fmt.Printf("\033[%d;1H", posY)
	if !first && tutStep[0] == 1 && tutStep[1] == 0 {
		r := bufio.NewReader(os.Stdin)
		Print("\033[1m<Rozpocznie się samouczek, jeśli chcesz go pominąć wpisz \"skip\", inaczej poprostu naciśnij enter>\033[0m")
		o, _ := r.ReadString('\n')
		fmt.Print("\033[A\033[2K")
		if o == "skip" {
			first = true
		}
	}
	if !first {
		Print("\r\033[A\033[2K")
		PrintLine("Samouczek")
		debug(tutStep)
		debug(formatStep(tutStep))
		switch formatStep(tutStep) {
		case "1,0":
			Talk([][2]interface{}{
				{"Mama", "Cześć, czy możesz mi pomóc?"},
				{"Ty", "Jasne, o co chodzi?"},
				{"Mama", "Dziękuje! jeśli możesz, skocz do sklepu i kup gazetę"},
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
			tip("Udaj się na dwór, tam znajduje się mały sklepik w którym możesz zakupić gazety")
			tutStep = []int{1, 1}
		case "1,1", "1,2", "1,3":
			Talk([][2]interface{}{
				{"Mama", "Halo, gdzie te zakupy?"},
				{"Ty", func() string {
					Print("Odpowiedz dlaczego nie kupiłeś gazety!!")
					PrintC("1. Tak, zaraz to zrobię, mamo...")
					PrintC("2. Przepraszam, ale nie mogę")
					for {
						inp := Prompt(">>> ")
						if inp == "1" {
							return "Tak, zaraz to zrobię, mamo..."
						} else if inp == "2" {
							return "Przepraszam, ale nie mogę"
						} else {
							Println("Niepoprawna opcja")
							continue
						}
					}
				}},
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
		case "2,0":
			time.Sleep(100 * ms)
			Tell("\033[1;33mMama\033[0;0m", Sprintf("%sDziękuje! a i jeszcze jedno, czy możesz wpłacić te pieniądze do banku?\033[0;0m", colorCodes["yellow"]))
			time.Sleep(1000 * ms)
			Println("\033[1;32m+ Dodano 200zł\033[0m")
			walletAdd(200)
			tip("Udaj się na dwór, tam znajduje się bank - w którym możesz wpłacić pieniądze")
		case "2,1", "2,2", "2,3":
			Talk([][2]interface{}{
				{"Mama", "Halo? dlaczego nie wpłaciłeś tych pierzonych pieniędzy??"},
				{"Ty", func() string {
					Print("Odpowiedz, dlaczego nie wykonałeś zadania!")
					PrintC("1. Tak, zaraz to zrobię, mamo...")
					PrintC("2. Przepraszam, ale nie mogę")
					for {
						inp := Prompt(">>> ")
						if inp == "1" {
							return "Tak, zaraz to zrobię, mamo..."
						} else if inp == "2" {
							return "Przepraszam, ale nie mogę"
						} else {
							ShowError("Niepoprawna opcja")
							continue
						}
					}
				}},
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
		case "3,0":
			Talk([][2]interface{}{
				{"Mama", "Dzięki, to już wszystko!"},
				{"Mama", "Teraz masz czas wolny"},
				{"Ty", "ok"}, // TODO: zmień to gówno
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
			first = true
		}
		Sep()
	}
	Println("Co chcesz teraz zrobić?")
	if !tmpses {
		commands("Wyjdź na dwór", "Konfiguracja", "Wyjdź z gry")
		for {
			inp := Prompt(">>> ")
			switch inp {
			case "1":
				goTo("DWÓR")
			case "2":
				goTo("CONFIG")
			case "3":
				Exit()
			default:
				ShowError("Nieznana opcja!")
				continue
			}
		}
	} else {
		commands("Wyjdź na dwór", "Wyjdź z gry")
		for {
			inp := Prompt(">>> ")
			switch inp {
			case "1":
				goTo("DWÓR")
			case "2":
				Exit()
			default:
				ShowError("Nieznana opcja!")
				continue
			}
		}
	}
}

func HandleOutside() {
	Println("Piękny zapach świeżego powietrza nie prawdasz?")
	for {
		commands("Udaj się do małego sklepu \"Ropucha\"", "Zagadaj do kogoś", "Zajdź do banku", "Idź do kasyna", "Powrót do domu")
		inp := Prompt(">>> ")
		switch inp {
		case "1":
			goTo("ROPUCHA")
		case "2":
			Talk(randTalk().([][2]interface{}), map[string]string{
				"Mama":     "yellow",
				"Obsługa":  "magenta",
				"Bezdomny": "brown",
				"Kupier":   "blue",
				"Ty":       "green",
			})
		case "3":
			goTo("BANK")
		case "4":
			goTo("KASYNO")
		case "5":
			goTo("DOM")
		default:
			PrintClr("Error!", "orange")
			PrintClr("Nieznana opcja", "red")
			continue
		}
	}
}
