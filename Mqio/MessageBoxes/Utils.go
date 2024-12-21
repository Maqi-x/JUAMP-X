package MessageBoxes

import (
	"regexp"
	"unicode/utf8"
)

func UnAnsii(input string) string {
	X := regexp.MustCompile(`\x1b\[[0-9;]*m`) // SHIT
	return X.ReplaceAllString(input, "")
}

func crunes(obj string) int {
    return utf8.RuneCountInString(obj)
}
