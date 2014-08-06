package main

import (
	mrand "math/rand"
	"strings"
	"time"

	"github.com/zerklabs/fifoq/string"
)

var (
	alphabet = []string{
		"",
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
		"u",
		"v",
		"w",
		"x",
		"y",
		"z",
	}

	characters = []string{
		"",
		"!",
		"*",
		"$",
		"%",
		"#",
		"^",
		"@",
	}

	numbers = []string{
		"",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"0",
	}

	av *fifoqstring.FifoqString
)

func generateStrongPassword(minLength int) string {
	av = fifoqstring.New(7)
	var password string

	mrand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < minLength; i++ {
		var lchar string

		for {
			set := mrand.Intn(4)

			if set == 1 {
				lchar = getAlphabet()
			} else if set == 2 {
				lchar = getNumber()
			} else if set == 3 {
				lchar = getSpecialCharacter()
			}

			if len(lchar) > 0 {
				if !av.Contains(lchar) {
					av.Push(lchar)
					break
				} else {
					mrand.Seed(time.Now().UTC().UnixNano())
				}
			}
		}

		password += lchar
	}

	return password
}

func getSpecialCharacter() string {
	pos := mrand.Intn(len(characters))

	return characters[pos]
}

func getNumber() string {
	pos := mrand.Intn(len(numbers))

	return numbers[pos]
}

func getAlphabet() string {
	pos := mrand.Intn(len(alphabet))

	toUpper := mrand.Intn(3)
	if toUpper == 0 || toUpper == 2 {
		return alphabet[pos]
	} else {
		return strings.ToUpper(alphabet[pos])
	}

	return alphabet[pos]
}
