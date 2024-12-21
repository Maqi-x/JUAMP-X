package main

import (
	mb "Mqio/MessageBoxes"
	"fmt"
	"strings"
	"time"
)

func main() {
	termCheck()
	con()
	chm("SAVE LOBBY")
	Println("Witaj! Juamp X to symulator życia, aby zacząć wybierz save (wpisz nazwę save, jeśli jeszcze nie mas save lub chcesz zacząć na nowym, to wpisz poprostu jaką chcesz nazwę, a zostanie utworzony nowy Save)")
	save = Prompt(">>> ")
	exitCode, err := loadSave(save)
	if err != nil {
		Println("wystąpił błąd podczas ładowania save:" + err.Error())
		Println("Zostaną wczytane domyślne wartości")
	}
	go func() {
		for {
			if !started {
				time.Sleep(10 * sec)
				continue
			}
			if autosave != -1 {
				saveSave(save)
				time.Sleep(time.Duration(autosave) * time.Millisecond)
			} else {
				time.Sleep(10 * sec)
				continue
			}
		}
	}()
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
		if w < 125 || h < 30 {
			m := mb.NewErrorBox(fmt.Sprintf("Aby uzyskać najlepsze wrażenia z gry prosimy zwiększyć rozmiar terminala, zapewni to lepszy wygląd oraz brak błędów graficznych. aktualnie ustawiłeś terminal na %d szerokości i %d wysokości, gdzie zalecany rozmiar to 125x30. Jeśli możesz zwiększ trochę terminal", w, h))
			m.Show()
			m.Hide()
			cleanT()
		} else if w < 125 {
			m := mb.NewErrorBox(fmt.Sprintf("Aby uzyskać najlepsze wrażenia z gry prosimy zwiększyć rozmiar terminala, zapewni to lepszy wygląd oraz brak błędów graficznych. aktualnie ustawiłeś terminal na %d szerokości i %d wysokości, gdzie zalecany rozmiar to 125x30. Jeśli możesz poszerz trochę terminal", w, h))
			m.Show()
			m.Hide()
			cleanT()
		} else if h < 30 {
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
	Println("Witaj w Juamp X!")
	Println("Najpierw podaj swoje imie...")
	NAME = Prompt(">>> ")
	Println(Sprintf("Okej %s, to teraz wylosujemy twoje miasto!", NAME))
	loading(1, "Losowanie")
	t := getTown()
	TOWN = t[0]
	Println("Wylosowałeś miasto: " + TOWN)
	for _, line := range strings.Split(t[1], "\n") {
		PrintC(line)
	}

	Println("")
	Println("A więc dobrze, zaczynajmy!")
	loading(4, "Podróż do lobby")
	saveSave(save)
	started = true
	goTo("LOBBY")
}
