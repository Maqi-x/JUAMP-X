package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func handleConfig() {
	w, _ := getTerminalSize()
	w = (w - lenr("Konfiguracja")) / 2
	Println(strings.Repeat("-", w) + "Konfiguracja" + strings.Repeat("-", w))
	for {
		PrintC("1. Zmienianie nazwy")
		PrintC("2. Ustawienia AutoSave")
		PrintC("3. Usuń save")
		PrintC("4. Zapisz ustawienia")
		PrintC("5. Powrót do domu")
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
			goTo("DOM")
		default:
			PrintClr("Error!", "orange")
			Println("Nieznana opcja")
		}
	}
}

func normTime(timeStr string) (int, error) {
	multiplier := 1

	switch {
	case strings.HasSuffix(timeStr, "ms"):
		timeStr = strings.TrimSuffix(timeStr, "ms")
		multiplier = 1
	case strings.HasSuffix(timeStr, "s"):
		timeStr = strings.TrimSuffix(timeStr, "s")
		multiplier = 1000
	case strings.HasSuffix(timeStr, "sec"):
		timeStr = strings.TrimSuffix(timeStr, "sec")
		multiplier = 1000
	case strings.HasSuffix(timeStr, "min"):
		timeStr = strings.TrimSuffix(timeStr, "min")
		multiplier = 60 * 1000
	case strings.HasSuffix(timeStr, "h"):
		timeStr = strings.TrimSuffix(timeStr, "h")
		multiplier = 60 * 60 * 1000
	default:
		return 0, errors.New("nie poprawny format czasu")
	}

	value, err := strconv.Atoi(strings.TrimSpace(timeStr))
	if err != nil {
		return 0, errors.New("Podany tekst to nie liczba")
	}

	return value * multiplier, nil
}
