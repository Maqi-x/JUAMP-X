package main

import (
	"os"
	"strings"
)

func handleLobby() {
	Println("Witaj! gdzie checsz sie udać?")
	PrintC("1. Ropucha")
	PrintC("2. Konfiguracja")
	PrintC("3. Wyjdź z gry")
	inp := Prompt(">>> ")
	switch inp {
	case "1":
		goTo("ROPUCHA")
	case "2":
		goTo("CONFIG")
	case "3":
		Exit()
	}
}

func handleConfig() {
	w, _ := getTerminalSize()
	w = (w - lenr("Konfiguracja")) / 2
	Println(strings.Repeat("-", w) + "Konfiguracja" + strings.Repeat("-", w))
	for {
		PrintC("1. Zmienianie nazwy")
		PrintC("2. Ustawienia AutoSave")
		PrintC("3. Usuń save")
		PrintC("4. Zapisz ustawienia")
		PrintC("5. Powrót do lobby")
		inp := Prompt(">>> ")
		Println("")
		switch inp {
		case "1":
			Println("Wpisz nową nazwe")
			NAME = Prompt(">>> ")
			PrintClr("Pomyślnie zapisano!", "green")
		case "2":
			func() {
				Println("Autosave: podaj czas co jaki AutoSave ma działać, np. 1h, 20m, 10s, 1m")
				PrintC("Jeśli chcesz wyłączyć AutoSave wpisz -1 lub 0")
				tm := Prompt(">>> ")
				if tm == "-1" || tm == "0" {
					autosave = -1
				} else {
					tmp, err := normTime(tm)
					if err != nil {
						PrintClr("Error!", "orange")
						PrintClr(err.Error(), "red")
						Println("")
						Println("ustawienia nie zostaną zastosowane")
						return
					}
					autosave = tmp
					PrintClr("Pomyślnie zapisano ustawienia!", "green")
				}
			}()
		case "3":
			err := os.Remove(Sprintf("saves/%s.toml", save))
			if err != nil {
				PrintClr("Error!", "orange")
				PrintClr("Błąd podczas usuwania save:"+err.Error(), "red")
			} else {
				PrintClr("Plik został pomyślnie usunięty", "green")
			}
		case "4":
			saveSave(save)
			PrintClr("Pomyślnie zapisano ustawienia!", "green")
		case "5":
			goTo("LOBBY")
		default:
			PrintClr("Error!", "orange")
			Println("Nieznana opcja")
		}
	}
}
