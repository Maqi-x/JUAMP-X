package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var RTitle bool
var posY int = 0

func chm(title string) {
	PLACE = title
	cleanT()
	renderTitle(PLACE)
}

func getTerminalSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0
	}
	return width, height
}

func clearT() {
	cleanT()
}

func cleanT() {
	fmt.Print("\033[2J\033[H")
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	posY = 1
}

// Title rendering
func renderTitle(title string) {
	width, height := getTerminalSize()
	border := strings.Repeat("#", width)
	padding := ((width - len(title)) - 2) / 2
	if padding < 0 {
		padding = 0
	}

	fTitle := fmt.Sprintf("%s %s %s", strings.Repeat("#", padding), title, strings.Repeat("#", padding))
	fmt.Printf("\033[%d;1H%s\n%s\n\033[%d;1H%s\n\033[%d;1H%s\n",
		height-3, fmt.Sprintf("Stan konta - portfel: %.2f, bank: %.2f; Głód: %d%s", wallet, bank, hungry, func() string {
			if hungry < 10 {
				return ", jesteś BARDZO GŁODNY"
			} else if hungry < 20 {
				return ", jesteś głodny"
			} else {
				return ""
			}
		}()), border,
		height-1, fTitle,
		height, border)
}

func Println(txt string) {
	width, height := getTerminalSize()

	lines := (lenr(txt) / width) + 1
	if !RTitle {
		cleanT()
		renderTitle(PLACE)
		RTitle = true
	}

	if posY+lines+2 > height {
		cleanT()
		renderTitle(PLACE)
	}

	fmt.Printf("\033[%d;1H", posY)
	fmt.Println(txt)
	posY += lines
}

func PrintClr(txt string, color string) {
	width, height := getTerminalSize()

	lines := (lenr(txt) / width) + 1
	if !RTitle {
		cleanT()
		renderTitle(PLACE)
		RTitle = true
	}

	if posY+lines+2 > height {
		cleanT()
		renderTitle(PLACE)
	}

	fmt.Printf("\033[%d;1H", posY)
	fmt.Println(colorCodes[color] + txt + "\033[0m")
	posY += lines
}

func Print(txt string) {
	fmt.Print(txt)
}

func PrintC(txt string) {
	width, height := getTerminalSize()

	lines := (lenr(txt) / width) + 1
	if !RTitle {
		cleanT()
		renderTitle(PLACE)
		RTitle = true
	}

	if posY+lines+3 > height {
		cleanT()
		renderTitle(PLACE)
	}

	// Text Centring
	padding := (width - len(txt)) / 2
	if padding < 0 {
		padding = 0
	}
	centeredText := fmt.Sprintf("%s%s", strings.Repeat(" ", padding), txt)

	fmt.Printf("\033[%d;1H%s\n", posY, centeredText)
	posY += lines
}

func Prompt(t string) string {
	setState(true)
	width, height := getTerminalSize()

	lines := (lenr(t) / width) + 1

	if !RTitle {
		cleanT()
		renderTitle(PLACE)
		RTitle = true
	}

	if posY+lines+4 > height {
		cleanT()
		renderTitle(PLACE)
	}

	fmt.Printf("\033[%d;1H%s", posY, t)
	posY += lines

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	posY += 1

	setState(false)
	return strings.TrimSpace(input)
}

func loading(secs int, prompt string) {
	duration := time.Duration(secs) * sec
	end := time.Now().Add(duration)

	for time.Now().Before(end) {
		for i := 1; i <= 3 && time.Now().Before(end); i++ {
			dots := strings.Repeat(".", i)
			fmt.Printf("\r%s%s", prompt, dots+strings.Repeat(" ", 3-i))
			time.Sleep(200 * ms)
		}
	}
	fmt.Print("\n")
}

func lenr(txt string) int {
	return utf8.RuneCountInString(txt)
}

func goTo(place string) {
	history = append(history, place)
	switch place {
	case "DOM":
		PLACE = "DOM"
		chm("DOM")
		handleLobby()
	case "ROPUCHA":
		PLACE = "ROPUCHA"
		chm("ROPUCHA")
		handleRopucha()
	case "CONFIG":
		PLACE = "CONFIG"
		chm("CONFIG")
		handleConfig()
	case "DWÓR":
		PLACE = "DWÓR"
		chm("DWÓR")
		HandleOutside()
	case "BANK":
		PLACE = "BANK"
		chm("BANK")
		handleBank()
	default:
		Println("Nieznana lokalizacja")
		Println("Wygląda na to, że save nie był poprawny! Nie powinno się ręcznie modefikować plików save!")
		loading(3, "Podróż do domu")
		goTo("DOM")
	}
}

func Back() {
	goTo(history[len(history)-2])
	history = history[:len(history)-1]
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func Exit(code ...int) {
	loading(1, "Zapisywanie save")
	err := saveSave(save)
	if err != nil {
		Println("Error podczas zapisywania save... przepraszamy za problemy, możliwe że save nie zostanie zapisany")
		time.Sleep(1 * sec)
	}
	clearT()
	if len(code) > 0 {
		os.Exit(code[0])
	}
	os.Exit(0)
}

func Tell(person string, txt string) {
	Println(Sprintf("[ %s ]: %s", person, txt))
}

func PrintLine(txt string) {
	txt = Sprintf(" %s ", txt)
	w, _ := getTerminalSize()
	w = (w - lenr(txt)) / 2
	Println(strings.Repeat("-", w) + txt + strings.Repeat("-", w))
}

func Sep() {
	w, _ := getTerminalSize()
	Println(strings.Repeat("-", w))
}

func normStr(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, str)
}

func hungryAdd(num int) {
	hungry += num
	if hungry > 100 {
		hungry = 100
	}
}

func Talk(msgs [][2]string, colors map[string]string) {
	for _, msg := range msgs {
		Println(Sprintf("\r\033[1m%s[ %s ]:\033[0m %s%s\033[0m",
			colorCodes[colors[msg[0]]], msg[0], colorCodes[colors[msg[0]]], msg[1],
		))
		r := bufio.NewReader(os.Stdin)
		Print("\033[30;47mDalej >\033[0m")
		r.ReadString('\n')
	}
}
