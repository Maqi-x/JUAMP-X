package main

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
	"strconv"
	"strings"
)

func autosaveSet(mode bool) (int, bool) {
	PrintColor("<bold><x>Ustawienia AutoSave</x></bold>")
	Println("Możesz tu wybrać co ile będzie działał autosave")
	Println("Autosave to system który automatycznie zapisuje stan gry, by zapobiec utratie postępu")
	if mode {
		PrintColor("Pamiętaj że wybierasz tylko ustawienia domyślne, czyli takie które będą automatycznie zastosowane w nowych save")
	}
	Println("Wpisz czas w formacie np. 10s, 10min, 10h, 10ms")
	PrintColor(Sprintf("<bold><x>Aktualny stan: %s</x></bold>", func() string {
		if mode {
			return msToTime(int(settings["autosavetime"].(int64)))
		} else {
			return msToTime(autosave)
		}
	}()))
	var x int
	tm := Prompt(">>> ")
	if tm == "-1" || tm == "0" {
		x = -1
	} else {
		tmp, err := normTime(tm)
		if err != nil {
			ShowError(err.Error())
			Println("")
			Println("ustawienia nie zostaną zastosowane")
			return 0, false
		}
		x = tmp
		PrintS("Pomyślnie zapisano ustawienia")
	}
	return x, true
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
		PrintC("5. Powrót do domu")
		inp := Prompt(">>> ")
		Println("")
		switch inp {
		case "1":
			Println("Wpisz nową nazwe")
			NAME = Prompt(">>> ")
			PrintS("Pomyślnie zapisano")
		case "2":
			tmp, x := autosaveSet(false)
			if x {
				autosave = tmp
			}
		case "3":
			err := os.Remove(Sprintf("saves/%s.toml", save))
			if err != nil {
				ShowError("Wystąpił błąd podczas usuwania save")
				debug(err.Error())
			} else {
				PrintS("Pomyślnie zapisano ustawienia")
			}
		case "4":
			saveSave(save)
			PrintS("Pomyślnie zapisano ustawienia")
		case "5":
			goTo("DOM")
		default:
			ShowError("Nieznana opcja")
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
		return 0, errors.New("podany tekst to nie liczba")
	}

	return value * multiplier, nil
}

func loadSettings() map[string]interface{} {
	var settings map[string]interface{}

	if _, err := toml.DecodeFile("global-settings.toml", &settings); err != nil {
		createDefault("global-settings.toml")
		return defaultSettings
	}

	for key, defaultValue := range defaultSettings {
		value, ok := settings[key]
		if !ok || !validateType(value, defaultValue) {
			Printfln("Nieprawidłowy lub brakujący klucz '%s'. Tworzę nowy plik z domyślnymi ustawieniami.", key)
			createDefault("global-settings.toml")
			return defaultSettings
		}
	}

	return settings
}

func validateType(value, defaultValue interface{}) bool {
	switch defaultValue.(type) {
	case int:
		_, x := value.(int64)
		return x
	case bool:
		_, x := value.(bool)
		return x
	case string:
		_, x := value.(string)
		return x
	default:
		return false
	}
}

func createDefault(fileName string) {
	_ = os.Remove(fileName)

	file, err := os.Create(fileName)
	if err != nil {
		Printfln("Error! nie można utworzyć pliku %s: %v", fileName, err)
		return
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(defaultSettings)
	if err != nil {
		ShowErrorf("nie można zapisać domyślnych ustawień do pliku %s: %v", fileName, err)
	}
}

func saveSettings() {
	file, err := os.Create("global-settings.toml")
	if err != nil {
		ShowError("Błąd: nie można utworzyć pliku: \n")
		debug(err.Error())
		return
	}
	defer file.Close()

	normalizedSettings := normSettings(settings)
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(normalizedSettings)
	if err != nil {
		ShowError("Nie można zapisać ustawień do pliku")
		debug(err.Error())
		return
	}
}

func normSettings(settings map[string]interface{}) map[string]interface{} {
	norm := make(map[string]interface{})
	for key, value := range settings {
		switch v := value.(type) {
		case int64:
			norm[key] = int(v)
		default:
			norm[key] = v
		}
	}
	return norm
}

func globalSettings() {
	PrintLine("Ustawienia globalne")
	Println("Witaj w globalnych ustawieniach")
	Println("Te ustawienia są niezależne od save, są globalne dla całej gry")
	Println("")
	info := func() {
		commands("Domyślny Autosave", "automatyczne pomijanie samouczka", "Kolory i style", "motyw przewodni", "Opcje developerskie", "Zastosuj i zapisz ustawienia", "Wyjście")
	}
	info()
	for {
		inp := Prompt(">>> ")
		switch inp {
		case "1":
			tmp, t := autosaveSet(true)
			if t {
				autosave = tmp
			}
		case "2":
			PrintColor("<bold><x>Automatyczne pomijanie samouczka</x></bold>")
			Println("Ta opcja sprawia, że przy każdym kolejnym utworzeniu nowego save")
			Println("Samouczek będzie pomijany")
			PrintColor(Sprintf("<bold><x>Aktualny stan:</x></bold> %s", func() string {
				if settings["tutskip"].(bool) {
					return "Tak"
				} else {
					return "Nie"
				}
			}()))
			Println("Czy chcesz włączyć tą opcje?")
			commands("Tak", "Nie")
			for {
				inp = Prompt(">>> ")
				if inp == "1" {
					settings["tutskip"] = true
					break
				} else if inp == "2" {
					settings["tutskip"] = false
					break
				} else {
					Println("Niepoprawna opcja")
					continue
				}
			}
		case "3":
			PrintColor("<bold><x>Kolory i style</x></bold>")
			Println("Ta opcja sprawia, że tekst w grze ma różne kolory i style")
			Println("Jeśli to włączysz większość kolorowania tesktu powinno zniknąć")
			Println("Jednak może to nie działać w każdym przypadku")
			PrintColor(Sprintf("<bold><x>Aktualny stan:</x></bold> %s", func() string {
				x := settings["colors"].(bool)
				if x {
					return "Tak"
				} else {
					return "Nie"
				}
			}()))
			Println("Czy chcesz włączyć tą opcje?")
			commands("Tak", "Nie")
			for {
				inp = Prompt(">>> ")
				if inp == "1" {
					settings["colors"] = true
					break
				} else if inp == "2" {
					settings["colors"] = false
					break
				} else {
					Println("Niepoprawna opcja")
					continue
				}
			}
		case "4":
			PrintColor("<bold><x>Motyw przewodni</x></bold>")
			PrintColor("Ta opcja pozwala na wybranie motywu kolorystycznego gry")
			PrintColor("<bold><x>Aktualny stan:</x></bold> " + settings["theme"].(string))
			Println("")
			Println("Podaj kolor np. red, blue, cyan, green, yellow, magenta, white, black...")
			for {
				clr := Prompt(">>> ")
				if in(clr, colorCodes) && clr != "bold" && clr != "dim" && clr != "italic" && clr != "hidden" {
					settings["theme"] = clr
					PrintS("Pomyślnie zastosowano motyw")
					break
				} else {
					PrintColor("<bold>Niestety ale nie można zastosować tego kolory</bold>, czy chcesz zobaczyć liste dostępnych kolorów?")
					Print("[t/n] ")
					inp := Prompt(">>> ")
					if inp == "t" || inp == "T" {
						for x, y := range colorCodes {
							if x == "bold" || x == "dim" || x == "italic" || x == "hidden" {
								continue
							}
							Printfln("%s%s\033[0m", y, x)
							dalej()
						}
					}
				}
			}
		case "5":
			PrintColor("<bold><x>Opcje developerskie</x></bold>")
			Println("Ta opcja włącza debugowanie i inne przydatne wiadomości dla developerów")
			Println("Nie zaleca się jej włączać jeśli nie jesteś deweloperem")
			PrintColor(Sprintf("<bold><x>Aktualny stan:</x></bold> %s", func() string {
				if settings["devoptions"].(bool) {
					return "Tak"
				} else {
					return "Nie"
				}
			}()))
			Println("Czy chcesz włączyć tą opcje?")
			commands("Tak", "Nie")
			for {
				inp = Prompt(">>> ")
				if inp == "1" {
					settings["devoptions"] = true
					break
				} else if inp == "2" {
					settings["devoptions"] = false
					break
				} else {
					Println("Niepoprawna opcja")
					continue
				}
			}
		case "6":
			saveSettings()
			useSettings()
		case "7":
			saveSettings()
			return
		default:
			PrintColor("<orange>Error</orange>! <red>Nieznana opcja<red>")
		}
		info()
	}
}
