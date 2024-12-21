package MessageBoxes

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

// ----------------------------------- DATA ---------------------------------- \\

var linesE []string
var currentLineE strings.Builder

var (
	resetE   = "\033[0m"
	boldE    = "\033[1m"
	redE     = "\033[91m"
	upE      = "\033[%dA" // Cursor up ANSII
	ClearlnE = "\033[2K"
)

type ErrorBox struct {
	wE         int
	margE      int
	TextE      string
	TextWrapE  []string
	LineCountE int
	selectedE  int
}

//-----------------------------------------------------------------------------\\

func NewErrorBox(txtE string) ErrorBox {
	wE, _, errE := term.GetSize(int(os.Stdout.Fd()))
	if errE != nil {
		wE = 80
	}
	if wE < 40 {
		wE = 40
	}
	margE := 4
	textwE := wE - 2*margE - 2
	if textwE < 1 {
		textwE = 1
	}
	linesE := wrapE(txtE, textwE)
	lineCountE := 2 + len(linesE)

	return ErrorBox{
		wE:         wE,
		margE:      margE,
		TextE:      txtE,
		TextWrapE:  linesE,
		LineCountE: lineCountE,
		selectedE:  0,
	}
}

func (self ErrorBox) Show() {
	fmt.Println(renderlnE(self.wE, "[ERROR]"))
	for _, lineE := range self.TextWrapE {
		fmt.Println(centerLineE(lineE, self.wE, self.margE))
	}
	fmt.Println(renderlnE(self.wE, ""))

	buttonE := "OK"

	renderButtonE := func() {
		buttonTextE := fmt.Sprintf("\033[48;5;15m\033[30m %s \033[0m", buttonE) // Active Button Highlighting ANSI
		cleanTextE := UnAnsii(buttonTextE)
		buttonWidthE := crunes(cleanTextE)
		paddE := (self.wE - buttonWidthE) / 2
		fmt.Printf("\r%s%s%s", strings.Repeat(" ", self.wE), "\r", strings.Repeat(" ", paddE)+buttonTextE)
	}

	oldStateE, errE := term.MakeRaw(int(os.Stdin.Fd()))
	if errE != nil {
		panic(errE)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldStateE)
	renderButtonE()

	bufE := make([]byte, 3)
	for {
		_, errE := os.Stdin.Read(bufE)
		if errE != nil {
			panic(errE)
		}

		if bufE[0] == 13 { // Enter
			break
		}
	}
}

func (self ErrorBox) Hide() {
	totalLinesE := self.LineCountE - 1 // Button
	fmt.Printf(upE, totalLinesE)
	for i := 0; i < totalLinesE; i++ {
		fmt.Print(ClearlnE + "\r")
		if i < totalLinesE-1 {
			fmt.Print("\033[1B") // Back 1 line
		}
	}

	fmt.Printf(upE, totalLinesE-1)
	fmt.Print(ClearlnE)
	fmt.Print("\033[1A")
	fmt.Print("\033[J")
}

func wrapE(txtE string, wE int) []string {
	wordsE := strings.Fields(txtE)

	for _, wordE := range wordsE {
		if currentLineE.Len()+crunes(wordE)+1 > wE {
			linesE = append(linesE, currentLineE.String())
			currentLineE.Reset()
		}
		if currentLineE.Len() > 0 {
			currentLineE.WriteRune(' ')
		}
		currentLineE.WriteString(wordE)
	}

	if currentLineE.Len() > 0 {
		linesE = append(linesE, currentLineE.String())
	}

	return linesE
}

func renderlnE(wE int, labelE string) string {
	rE := func(oE string, iE int) string {
		return strings.Repeat(oE, iE)
	}
	labelE = " " + labelE + " "
	labelLenE := crunes(labelE)

	if labelLenE > wE {
		labelE = labelE[:wE]
		labelLenE = wE
	}

	paddE := (wE - labelLenE) / 2
	lineE := rE("-", paddE) + labelE + rE("-", wE-paddE-labelLenE)
	lineE = lineE[:wE-2]

	return fmt.Sprintf("%s|%s|%s",
		redE, lineE, resetE)
}

func centerLineE(txtE string, wE int, margE int) string {
	textLengthE := crunes(txtE)
	paddE := wE - textLengthE - 2*margE - 2
	if paddE < 0 {
		paddE = 0
	}
	leftpaddE := margE + paddE/2
	rightpaddE := margE + (paddE+1)/2
	return fmt.Sprintf("%s|%s%s%s|%s",
		redE, resetE,
		strings.Repeat(" ", leftpaddE)+txtE+strings.Repeat(" ", rightpaddE),
		redE, resetE)
}
