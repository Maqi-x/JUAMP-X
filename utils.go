package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"
)

var RTitle bool
var posY int = 0

func chm(title string) {
	PLACE = " " + title + " "
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
	exec.Command("clear")
	posY = 1
}

// Title rendering
func renderTitle(title string) {
	width, height := getTerminalSize()
	border := strings.Repeat("#", width)
	padding := (width - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	fTitle := fmt.Sprintf("%s%s%s", strings.Repeat("#", padding), title, strings.Repeat("#", padding))
	fmt.Printf("\033[%d;1H%s\n%s\n\033[%d;1H%s\n\033[%d;1H%s\n",
		height-3, fmt.Sprintf("Stan konta - portfel: %d, bank: %d", wallet, bank), border,
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
	switch place {
	case "LOBBY":
		PLACE = "LOBBY"
		chm("LOBBY")
		handleLobby()
	case "ROPUCHA":
		PLACE = "ROPUCHA"
		chm("ROPUCHA")
		handleRopucha()
	}
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
