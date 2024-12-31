package main

import (
	mb "Mqio/MessageBoxes"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var exitCode int

func main() {
	defer func() {
		if r := recover(); r != nil {
			cleanT()
			debugf("panic: %v\n", r)
			fmt.Printf("%s%sWygląda na to że wystąpił błąd\033[0m przepraszamy za wszelkie niedogodności\n", colorCodes["red"], colorCodes["bold"])
			fmt.Println("Jednak nie ma co się denerwować - save zostanie zapisany tak samo jak i ustawienia")
			fmt.Println("\033[0mProsimy o zgłoszenie błędu do developerów\033[0m")
			dalej()
			saveSettings()
			saveSave(save)
			if devmode || in("--debug", os.Args) || in("--d", os.Args) {
				fmt.Printf("Czy chcesz zobaczyć traceback? (T/N) ")
				r := bufio.NewReader(os.Stdin)
				inp := strings.TrimSpace(func() string {
					inp, _ := r.ReadString('\n')
					return inp
				}())
				debug(inp)
				if inp == "T" || inp == "t" {
					panic(r)
				}
			}
		}
	}()
	debugf(strconv.Itoa((os.Getpid())))
	if devmode {
		dalej()
	}
	colorsCodes()
	termCheck()
	con()
	settings = loadSettings()
	useSettings = func() {
		if v, ok := settings["autosavetime"].(int64); ok {
			autosave = int(v)
		} else if v, ok := settings["autosavetime"].(int); ok {
			autosave = v
		}
		devmode = settings["devoptions"].(bool)
		for _, arg := range os.Args {
			if arg == "--debug" || arg == "-d" {
				devmode = true
			}
		}
		if !settings["colors"].(bool) {
			for key, _ := range colorCodes {
				colorCodes[key] = ""
			}
		}
		colorCodes["x"] = colorCodes[settings["theme"].(string)]
	}
	useSettings()
	chm("SAVE LOBBY")
	info := func() {
		PrintColor("<x><bold>Witaj!</bold></x> <x>Juamp X to symulator życia</x>, aby zacząć wybierz save (wpisz nazwę save, jeśli jeszcze nie mas save lub chcesz zacząć na nowym, to wpisz poprostu jaką chcesz nazwę, a zostanie utworzony nowy Save)")
		Println()
		PrintColor("Możesz też wpisać <bold>/settings</bold> by przejść do ustawień globalnych lub <bold>/exit</bold> by wyjść z gry")
		PrintColor("Dodatkowo <bold>/help</bold> pozwoli Ci na uzyskanie pomocy, a <bold>/list</bold> wyświetli liste save'ów")
		PrintColor("<bold>/saves</bold> otworzy Save Menager, tam możesz zarządzać swoimi save'ami")
		PrintColor("Jeśli chcesz coś wypróbować możesz użyć gry tymczasowej dzięki <bold>/tmp</bold>")
	}
	info()
	for {
		inp := Prompt(">>> ")
		if inp == "/exit" {
			os.Exit(0)
		} else if inp == "/settings" {
			goTo("GLOBALSETTINGS")
			chm("SAVE LOBBY")
			info()
			continue
		} else if inp == "/help" {
			info()
			continue
		} else if inp == "/list" {
			listSaves()
			continue
		} else if inp == "/tmp" {
			tmpses = true
			initData()
			goTo("DOM")
		} else if inp == "/saves" {
			PLACE = "SAVEMENAGER"
			chm("SAVE MENAGER")
			tmp, x := SaveMenager()
			if x {
				save = tmp
				var err error
				exitCode, err = loadSave(save)
				if err != nil {
					Println("wystąpił błąd podczas ładowania save:" + err.Error())
					Println("Zostaną wczytane domyślne wartości")
				}
				break
			}
			chm("SAVE LOBBY")
			info()
			continue
		} else if rg := regexp.MustCompile(`^[^/\\]*$`); rg.MatchString(inp) && strings.TrimSpace(inp) != "" {
			var err error
			save = inp
			exitCode, err = loadSave(inp)
			if err != nil {
				Println("wystąpił błąd podczas ładowania save:" + err.Error())
				Println("Zostaną wczytane domyślne wartości")
			}
			break
		}
		Println("Nie poprawny format, wpisz ponownie")
		Println("")
	}
	// ----------------------------------------------------- AutoSave --------------------------------------------------------- \\
	go func() {
		for {
			if !started {
				time.Sleep(10 * sec)
				continue
			}
			if autosave > 0 {
				saveSave(save)
				saveSettings()
				time.Sleep(time.Duration(autosave) * time.Millisecond)
			} else {
				time.Sleep(10 * sec)
				continue
			}
		}
	}()
	// ----------------------------------------------------- HungrySystem --------------------------------------------------------- \\
	go func() {
		var o, o1 bool
		for {
			if !started {
				time.Sleep(10 * sec)
				continue
			}
			if !(hungry < 0) {
				hungry--
			}
			time.Sleep(4500 * ms)
			if hungry < 20 && !o {
				PrintClr("Uwaga!", "orange")
				PrintClr("Jesteś głodny! udaj się jak najszybciej coś zjeść", "red")
				o = true
			} else if hungry < 10 && !o1 {
				PrintClr("UWAGA!", "orange")
				PrintClr("Jesteś BARDZO GŁODNY! udaj się jak najszybciej coś zjeść!", "red")
				o1 = true
			} else if hungry <= 0 {
				started = false
				err := os.Remove(Sprintf("saves/%s.toml", save))
				if err != nil {
					PrintClr("Error!", "orange")
					PrintClr(err.Error(), "red")
				}
				clearT()
				PLACE = "END"
				chm(PLACE)
				PrintC(colorCodes["red"] + "Koniec..." + colorCodes["reset"])
				PrintC(colorCodes["orange"] + "Dziękujemy za gre w JUAMP X!")
				PrintC("niestety, z powodu głodu nie przerzyłeś")
				PrintC("Jeśli chcesz kontynuować przygodę utwórz nowego save" + colorCodes["reset"])
				loading(5, "Koniec gry")
				os.Exit(0)
			}
		}
	}()
	go func() {
		for {
			time.Sleep(600 * sec)
			assortment["buns"] += randint(7, 11)
			assortment["pizzas"] += randint(5, 8)
			assortment["newspapers"] += randint(3, 7)
			assortment["shoes"] += randint(2, 4)
		}
	}()
	go func() {
		for {
			time.Sleep(180 * sec)
			if !started {
				continue
			}
			age += 1
		}
	}()
	// --------------------------------------------------------------------------------------------------------------------------- \\
	switch exitCode {
	case 0:
		Println("Pomyślnie załadowano save!")
		started = true
		goTo(PLACE)
	case 1:
		Println("Pomyślnie utworzono nowy save!")
		start()
	case 2:
		Println("Wygląda na to, że save nie był poprawny! Nie powinno się ręcznie modefikować plików save, plik zostanie nadpisany domyślnymi danymi")
		start()
	}
	clearT()
}

func con() {
	go func() {
		var x, y, x1, y1 int
		x, y = getTerminalSize()
		x1, y1 = getTerminalSize()
		for {
			if getState() {
				continue
			}
			x, y = getTerminalSize()
			if x != x1 || y != y1 {
				cleanT()
				renderTitle(PLACE)
				RTitle = true
			}
			x1, y1 = x, y
			time.Sleep(5 * ms)
		}
	}()
}

func termCheck() {
	cleanT()
	for {
		w, h := getTerminalSize()
		if w == 0 && h == 0 {
			m := mb.NewErrorBox("Wygląda na to że wczytanie rozmiaru terminala się nie powiodło, czy aby napewno chcesz kontynuować? mogą występować różne problemu graficzne, zalecamy skorzystać z innego terminala jeśli to rozwiązuje sproblem. jeśli tak naciśnij \"Ok\"")
			m.Show()
			m.Hide()
			clearT()
			return
		}
		if w < 145 && h < 35 {
			m := mb.NewErrorBox(fmt.Sprintf("Aby uzyskać najlepsze wrażenia z gry prosimy zwiększyć rozmiar terminala, zapewni to lepszy wygląd oraz brak błędów graficznych. aktualnie ustawiłeś terminal na %d szerokości i %d wysokości, gdzie zalecany rozmiar to 125x30. Jeśli możesz zwiększ trochę terminal", w, h))
			m.Show()
			m.Hide()
			cleanT()
		} else if w < 145 {
			m := mb.NewErrorBox(fmt.Sprintf("Aby uzyskać najlepsze wrażenia z gry prosimy zwiększyć rozmiar terminala, zapewni to lepszy wygląd oraz brak błędów graficznych. aktualnie ustawiłeś terminal na %d szerokości i %d wysokości, gdzie zalecany rozmiar to 125x30. Jeśli możesz poszerz trochę terminal", w, h))
			m.Show()
			m.Hide()
			cleanT()
		} else if h < 35 {
			m := mb.NewErrorBox(fmt.Sprintf("Aby uzyskać najlepsze wrażenia z gry prosimy zwiększyć rozmiar terminala, zapewni to lepszy wygląd oraz brak błędów graficznych. aktualnie ustawiłeś terminal na %d szerokości i %d wysokości, gdzie zalecany rozmiar to 125x30. Jeśli możesz zwyż trochę terminal", w, h))
			m.Show()
			m.Hide()
			cleanT()
		} else {
			cleanT()
			return
		}
		time.Sleep(5 * ms)
	}
}

func start() {
	debug(save)
	PrintColor("Witaj w <bold>Juamp X!</bold>")
	Println()
	dalej()
	PrintS("Nazwa/imie:")
	PrintColor("Aby zapewnić lepsze doświadczenie z gry podaj nam swój <bold>nick/imie<bold>")
	Println("Będzie ono używane głównie w dialogach i nie będzie miało wpływu na rozgrywke")
	NAME = Prompt(">>> ")
	Println()
	PrintColor(Sprintf("Okej <bold>%s</bold>, to teraz wylosujemy twoje <bold>miasto!</bold>", NAME))
	loading(1, "Losowanie")
	t := getTown()
	TOWN = t[0]
	Println("Wylosowałeś miasto: " + TOWN)
	inf := centerAlign(strings.Split(t[1], "\n"), func() int {
		w, _ := getTerminalSize()
		return w
	}())
	for _, line := range strings.Split(inf, "\n") {
		Println(line)
	}

	Println("")
	time.Sleep(1 * sec)
	dalej()
	Println("A więc dobrze, zaczynajmy!")
	loading(2, "Podróż do domu")
	saveSave(save)
	started = true
	first = settings["tutskip"].(bool)
	goTo("DOM")
}
