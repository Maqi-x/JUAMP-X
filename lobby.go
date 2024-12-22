package main

import (
	"fmt"
	"strings"
)

func formatStep(step []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(step), " ", ",", -1), "[]")
}

func handleLobby() {
	o := Prompt("Rozpocznie się samouczek, jeśli chcesz go pominąć wpisz \"skip\", inaczej poprostu naciśnij enter")
	if o == "skip" {
		first = true
	}
	if !first {
		switch formatStep(tutStep) {
		case "1,0":
			PrintLine("Samouczek")
			Talk([][2]string{
				{"Mama", "Cześć, czy możesz mi pomóc?"},
				{"Ty", "Jasne, o co chodzi?"},
				{"Mama", "Dziękuje! jeśli możesz, skocz do sklepu i kup gazetę"},
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
			PrintClr("\033[1mPorada:", "blue")
			PrintClr("Udaj się na dwór, tam znajduje się mały sklepik w którym możesz zakupić gazety", "blue")
			tutStep = []int{1, 1}
		case "1,1", "1,2", "1,3":
			Talk([][2]string{
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
				}()},
			}, map[string]string{
				"Mama": "yellow",
				"Ty":   "green",
			})
		case "2,0":
			Tell("\033[1;33mMama\033[0;0m", Sprintf("%sDziękuje! a i jeszcze jedno, czy możesz wpłacić te pieniądze do banku?\033[0;0m", colorCodes["yellow"]))
			PrintC("\033[1;32mDodano 100zł\033[0m   ")
			bankAdd(100)
			PrintClr("\033[1mPorada:", "blue")
			PrintClr("Udaj się na dwór, tam znajduje się bank - w którym możesz wpłacić pieniądze", "blue")
		}
	}
	Println("Co chcesz teraz zrobić?")
	PrintC("1. wyjdź na dwór")
	PrintC("2. Konfiguracja")
	PrintC("3. Wyjdź z gry")
	inp := Prompt(">>> ")
	switch inp {
	case "1":
		goTo("DWÓR")
	case "2":
		goTo("CONFIG")
	case "3":
		Exit()
	default:
		PrintClr("Error!", "orange")
		PrintClr("Nieznana opcja", "red")
	}
}

func HandleOutside() {
	Println("Piękny zapach świeżego powietrza nie prawdasz?")
	for {
		PrintC("1. Udaj się do małego sklepu \"Ropucha\"")
		PrintC("2. Zagadaj do kogoś")
		PrintC("3. Zajdź do banku")
		PrintC("4. Powrót do domu")
		inp := Prompt(">>> ")
		switch inp {
		case "1":
			goTo("ROPUCHA")
		case "2":
			// not inplemented
		case "3":
			goTo("BANK")
		case "4":
			goTo("DOM")
		default:
			PrintClr("Error!", "orange")
			PrintClr("Nieznana opcja", "red")
			continue
		}
	}
}
