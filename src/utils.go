package main

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
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
	fmt.Printf("\033[%d;1H%s\n\033[1m%s\n\033[%d;1H%s\n\033[%d;1H%s\n\033[0m",
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
	fmt.Print("\033[0m\033[H")
}

func Println(texts ...any) {
	var txt string
	for _, t := range texts {
		txt += Sprintf("%v", t) + " "
	}
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
	if t == ">>> " && settings["colors"].(bool) {
		t = "\033[1;32m>>>\033[0m \033[38;2;180;255;200m"
	}
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
	Print("\033[0m")
	return strings.TrimSpace(input)
}

func loading(secs int, prompt string) {
	duration := time.Duration(secs) * sec
	end := time.Now().Add(duration)

	for time.Now().Before(end) {
		for i := 1; i <= 3 && time.Now().Before(end); i++ {
			dots := strings.Repeat(".", i)
			Printf("\r%s%s", prompt, dots+strings.Repeat(" ", 3-i))
			time.Sleep(200 * ms)
		}
	}
	Println()
}

func lenr(txt string) int {
	ansiRegex := regexp.MustCompile(`\033\[[0-9;]*[a-zA-Z]`)

	cleanTxt := ansiRegex.ReplaceAllString(txt, "")

	return utf8.RuneCountInString(cleanTxt)
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
	case "KASYNO":
		PLACE = "KASYNO"
		chm("KASYNO")
		handleKasyno()
	case "GLOBALSETTINGS":
		PLACE = "GLOBAL SETTINGS"
		chm("GLOBAL SETTINGS")
		globalSettings()
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
	saveSettings()
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

func Talk(msgs [][2]interface{}, colors map[string]string) {
	for _, msg := range msgs {
		message := func() string {
			if fn, ok := msg[1].(func() string); ok {
				return fn()
			}
			if str, ok := msg[1].(string); ok {
				return str
			}
			return ""
		}()

		Println(Sprintf("\r\033[1m%s[ %s ]:\033[0m %s%s\033[0m",
			colorCodes[colors[msg[0].(string)]], msg[0], colorCodes[colors[msg[0].(string)]], message,
		))
		r := bufio.NewReader(os.Stdin)
		Print("\033[30;47mDalej >\033[0m")
		r.ReadString('\n')
		Print("\r")
	}
}

func centerAlign(lines []string, terminalWidth int) string {
	var result strings.Builder

	maxWidth := 0
	for _, line := range lines {
		if lenr(line) > maxWidth {
			maxWidth = lenr(line)
		}
	}

	for _, line := range lines {
		padding := (terminalWidth - maxWidth) / 2
		result.WriteString(strings.Repeat(" ", padding) + line + "\n")
	}

	return result.String()
}

func commands(cmds ...string) {
	if settings["colors"].(bool) {
		lines := []string{Sprintf("%sKomendy:\x1b[0m", colorCodes["x"]+"\x1b[1m")}
		for i, cmd := range cmds {
			lines = append(lines, fmt.Sprintf("%s%d.%s %s\033[0m", colorCodes["x"]+"\033[1m", i+1, "\033[0m"+colorCodes["x"], cmd))
		}
		for _, line := range strings.Split(centerAlign(lines, func() int {
			width, _ := getTerminalSize()
			return width - 1
		}()), "\n") {
			Println(line)
		}
		return
	}
	lines := []string{"Komendy:"}
	for i, cmd := range cmds {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, cmd))
	}
	for _, line := range strings.Split(centerAlign(lines, func() int {
		width, _ := getTerminalSize()
		return width - 1
	}()), "\n") {
		Println(line)
	}
	return
}

func listSaves() {
	var cmds []string
	err := filepath.WalkDir("saves", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".toml") {
			cmds = append(cmds, strings.TrimSuffix(d.Name(), ".toml"))
		}
		return nil
	})
	if err != nil {
		debugf("error while listing saves: %v\n", err)
	}
	if len(cmds) == 0 {
		if settings["colors"].(bool) {
			PrintS(centerAlign([]string{"Brak zapisanych gier"}, func() int {
				w, _ := getTerminalSize()
				return w - 1
			}()))
			return
		}
		PrintC("Brak zapisanych gier")
		return
	}
	if settings["colors"].(bool) {
		lines := []string{Sprintf("%sZapisane gry:\x1b[0m", colorCodes["x"]+"\x1b[1m")}
		for i, cmd := range cmds {
			lines = append(lines, fmt.Sprintf("%s%d.%s %s\033[0m", colorCodes["x"]+"\033[1m", i+1, "\033[0m"+colorCodes["x"], cmd))
		}
		for _, line := range strings.Split(centerAlign(lines, func() int {
			width, _ := getTerminalSize()
			return width - 1
		}()), "\n") {
			Println(line)
		}
		return
	}
	lines := []string{"Zapisane gry:"}
	for i, cmd := range cmds {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, cmd))
	}
	for _, line := range strings.Split(centerAlign(lines, func() int {
		width, _ := getTerminalSize()
		return width - 1
	}()), "\n") {
		Println(line)
	}
	return
}

func randint[ints interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}](min ints, max ints) int {
	minInt := int(min)
	maxInt := int(max)

	if minInt > maxInt {
		minInt, maxInt = maxInt, minInt
	}
	return rand.Intn(maxInt-minInt+1) + minInt
}

func Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func Printfln(format string, a ...interface{}) {
	Println(Sprintf(format, a...))
}

func dalej() {
	r := bufio.NewReader(os.Stdin)
	Print("\033[30;47mDalej >\033[0m")
	r.ReadString('\n')
	Print("\r")
}

func debugf(format string, args ...interface{}) {
	if devmode {
		debug(Sprintf(format, args...))
	}
}

func debug(msgs ...interface{}) {
	if devmode {
		var txt string
		for _, x := range msgs {
			txt += Sprintf("%v ", x)
		}
		Printfln("[DEBUG]: %s", txt)
	}
}

func PrintColor(txts ...any) {
	txt := Sprintf("%v", txts[0])
	for _, t := range txts[1:] {
		txt += Sprintf(" %v", t)
	}
	re := regexp.MustCompile(`<(\/?)([a-zA-Z_]+)>`)

	var Stack []string

	Println(re.ReplaceAllStringFunc(txt, func(match string) string {
		isClosing := match[1] == '/'
		tag := strings.Trim(match, "</>")

		if isClosing {
			for i := len(Stack) - 1; i >= 0; i-- {
				if Stack[i] == tag {
					Stack = append(Stack[:i], Stack[i+1:]...)
					break
				}
			}
			if len(Stack) > 0 {
				var activeStyles string
				for _, style := range Stack {
					if code, exists := colorCodes[style]; exists {
						activeStyles += code
					}
				}
				return activeStyles
			}
			return colorCodes["reset"]
		}
		if code, exists := colorCodes[tag]; exists {
			Stack = append(Stack, tag)
			return code
		}
		return match
	}) + "\033[0m")
}

func msToTime(milliseconds int) string {
	seconds := milliseconds / 1000
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	seconds %= 60
	minutes %= 60
	hours %= 24

	var result string

	if days > 0 {
		result += strconv.Itoa(days) + "d "
	}
	if hours > 0 {
		result += strconv.Itoa(hours) + "godz "
	}
	if minutes > 0 {
		result += strconv.Itoa(minutes) + "min "
	}
	if seconds > 0 {
		result += strconv.Itoa(seconds) + "sec "
	}

	if result == "" {
		result = "0sek"
	}

	return result
}

func in(element any, block any) bool {
	blockValue := reflect.ValueOf(block)

	if blockValue.Kind() == reflect.Map {
		return blockValue.MapIndex(reflect.ValueOf(element)).IsValid()
	}

	if blockValue.Kind() == reflect.Array || blockValue.Kind() == reflect.Slice {
		for i := 0; i < blockValue.Len(); i++ {
			if reflect.DeepEqual(blockValue.Index(i).Interface(), element) {
				return true
			}
		}
	}

	return false
}

func PrintS(txt string) {
	PrintColor(Sprintf("<x><bold>%s</bold></x>", txt))
}

func ShowError(errs ...any) {
	var str string
	for _, err := range errs {
		str += Sprintf(" %v", err)
	}
	if settings["colors"].(bool) {
		Printfln("\033[1;31mBłąd!\033[0;91m%s\033[0m", str)
	} else {
		Printfln("Błąd!%s", str)
	}
}

func ShowErrorf(format string, x ...any) {
	str := Sprintf(format, x...)
	if settings["colors"].(bool) {
		Printfln("\033[1;31mBłąd!\033[0;91m %s\033[0m", str)
	} else {
		Printfln("Błąd! %s", str)
	}
}

func tip(txts ...string) {
	if len(txts) == 0 {
		return
	}
	txt := Sprintf("%s", txts[0])
	for _, t := range txts[1:] {
		txt += Sprintf(" %v", t)
	}
	Printfln("\033[1m%sPorada:", colorCodes["bright_cyan"])
	Printfln("\033[22m%s\033[0m", txt)
}
