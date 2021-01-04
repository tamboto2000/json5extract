package json5extract

import (
	"unicode/utf8"
)

// Convert []rune to []byte
func runesToUTF8(rs []rune) []byte {
	size := 0
	for _, r := range rs {
		size += utf8.RuneLen(r)
	}

	bs := make([]byte, size)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}

	return bs
}

// Check if a rune is a hexa char
func isCharHex(char rune) bool {
	for _, c := range hex {
		if char == c {
			return true
		}
	}

	return false
}

func removeSliceIndex(s []rune, i int) []rune {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func isReservedWord(chars []rune) bool {
	rl := len(chars)
	for _, word := range rsvWords {
		if len(word) == rl {
			for i, c := range word {
				if chars[i] != c {
					return false
				}
			}

			return true
		}
	}

	return false
}
