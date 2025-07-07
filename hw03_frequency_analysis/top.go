package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordPair struct {
	word   string
	number int
}

func Top10(ourString string) []string {
	if strings.EqualFold(ourString, "") {
		return []string{}
	}
	fields := strings.Fields(ourString)

	mapString := make(map[string]int)
	for _, word := range fields {
		if strings.EqualFold(word, " ") || strings.EqualFold(word, "") {
			continue
		}
		_, ok := mapString[word]
		if ok {
			mapString[word]++
		} else {
			mapString[word] = 1
		}
	}

	sliceWord := make([]WordPair, 0)
	for key, value := range mapString {
		sliceWord = append(sliceWord, WordPair{key, value})
	}

	slice := sort.Slice
	slice(sliceWord, func(i, j int) bool {
		if sliceWord[i].number != sliceWord[j].number {
			return sliceWord[i].number > sliceWord[j].number
		}
		return sliceWord[i].word < sliceWord[j].word
	})

	resultLen := 10
	if len(sliceWord) < 10 {
		resultLen = len(sliceWord)
	}
	result := make([]string, resultLen)
	for i := 0; i < resultLen; i++ {
		result[i] = sliceWord[i].word
	}
	return result
}
