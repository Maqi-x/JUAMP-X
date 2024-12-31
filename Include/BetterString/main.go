package BetterString

import (
	"fmt"
	"regexp"
	"strings"
)

type BetterString string

func (s BetterString) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

func (s BetterString) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

func (s BetterString) TrimPrefix(prefix string) BetterString {
	return BetterString(strings.TrimPrefix(string(s), prefix))
}

func (s BetterString) TrimSuffix(suffix string) BetterString {
	return BetterString(strings.TrimSuffix(string(s), suffix))
}

func (s BetterString) TrimSpace() BetterString {
	return BetterString(strings.TrimSpace(string(s)))
}

func (s BetterString) ToUpper() BetterString {
	return BetterString(strings.ToUpper(string(s)))
}

func (s BetterString) ToLower() BetterString {
	return BetterString(strings.ToLower(string(s)))
}

func (s BetterString) Replace(old, new string) BetterString {
	return BetterString(strings.ReplaceAll(string(s), old, new))
}

func (s BetterString) Remove(substr string) BetterString {
	return BetterString(strings.ReplaceAll(string(s), substr, ""))
}

func (s BetterString) Repeat(count int) BetterString {
	return BetterString(strings.Repeat(string(s), count))
}

func (s BetterString) Add(other any) BetterString {
	switch v := other.(type) {
	case BetterString:
		return BetterString(string(s) + string(v))
	case string:
		return BetterString(string(s) + v)
	default:
		panic("unsuported type")
	}
}

func New(s string) BetterString {
	return BetterString(s)
}

func (s BetterString) Split(sep string) []BetterString {
	parts := strings.Split(string(s), sep)
	result := make([]BetterString, len(parts))
	for i, part := range parts {
		result[i] = BetterString(part)
	}
	return result
}

func (s BetterString) String() string {
	return string(s)
}

func (s BetterString) SplitString() ([]string, error) {
	input := string(s)
	re := regexp.MustCompile(`"([^"]*)"|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return nil, fmt.Errorf("brak dopasowań w ciągu: %s", input)
	}

	var result []string
	for _, match := range matches {
		if match[1] != "" {
			result = append(result, match[1])
		} else {
			result = append(result, match[2])
		}
	}
	return result, nil
}
