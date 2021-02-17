package json5extract

import (
	"bufio"
	"bytes"
	"strconv"
	"unicode"
)

var (
	hex = []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F'}
)

// String types
const (
	doubleQuotedStr = iota
	singleQuotedStr
)

func parseStr(r reader, ty int) (*JSON5, error) {
	str := &JSON5{kind: String}
	if ty == doubleQuotedStr {
		str.push('"')
	} else {
		str.push('\'')
	}

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		// detect escaped char
		if char == '\\' {
			str.push(char)

			char, _, err := r.ReadRune()
			if err != nil {
				return nil, err
			}

			// unicode
			if char == 'u' {
				str.push(char)
				for i := 0; i < 4; i++ {
					char, _, err := r.ReadRune()
					if err != nil {
						return nil, err
					}

					if !isCharHex(char) {
						return nil, ErrInvalidFormat
					}

					str.push(char)
				}

				continue
			}

			// hexa
			if char == 'x' {
				str.push(char)
				for i := 0; i < 2; i++ {
					char, _, err := r.ReadRune()
					if err != nil {
						return nil, err
					}

					if !isCharHex(char) {
						return nil, ErrInvalidFormat
					}

					str.push(char)
				}

				continue
			}

			// numeric
			if unicode.IsNumber(char) && char != '0' {
				return nil, ErrInvalidFormat
			}

			str.push(char)
			continue
		}

		// detect line terminator (line feed or return carriage and line feed)
		if char == '\n' {
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char == '\\' {
				prev2char := str.raw[rawlen-2]
				if prev2char == '\\' {
					return nil, ErrInvalidFormat
				}
			}

			if prev1char == '\r' {
				prev2char := str.raw[rawlen-2]
				if prev2char != '\\' {
					return nil, ErrInvalidFormat
				}

				prev3char := str.raw[rawlen-3]
				if prev3char == '\\' {
					return nil, ErrInvalidFormat
				}
			}

			str.push(char)
			continue
		}

		// detect line terminator (return carriage)
		if char == '\r' {
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char != '\\' {
				return nil, ErrInvalidFormat
			}

			prev2char := str.raw[rawlen-2]
			if prev2char == '\\' {
				return nil, ErrInvalidFormat
			}

			str.push(char)
			continue
		}

		// detect line terminator (line separator)
		if char == '\u2028' {
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char != '\\' {
				return nil, ErrInvalidFormat
			}

			prev2char := str.raw[rawlen-2]
			if prev2char == '\\' {
				return nil, ErrInvalidFormat
			}

			str.push(char)
			continue
		}

		// detect line terminator (line paragraph separator)
		if char == '\u2029' {
			str.push(char)
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char != '\\' {
				return nil, ErrInvalidFormat
			}

			prev2char := str.raw[rawlen-2]
			if prev2char == '\\' {
				return nil, ErrInvalidFormat
			}

			continue
		}

		// detect string punctuation (double quote)
		if char == '"' && ty == doubleQuotedStr {
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char != '\\' {
				str.push(char)
				break
			}

			prev2char := str.raw[rawlen-2]
			if prev2char == '\\' {
				str.push(char)
				break
			}

			str.push(char)
			continue
		}

		// detect string punctuation (single quote)
		if char == '\'' && ty == singleQuotedStr {
			rawlen := len(str.raw)
			prev1char := str.raw[rawlen-1]
			if prev1char != '\\' {
				str.push(char)
				break
			}

			prev2char := str.raw[rawlen-2]
			if prev2char == '\\' {
				str.push(char)
				break
			}

			str.push(char)
			continue
		}

		str.push(char)
	}

	str.val = unescapeRunesToStr(str.raw, ty)

	return str, nil
}

func unescapeRunesToStr(chars []rune, ty int) string {
	newRunes := make([]rune, 0)
	r := bufio.NewReader(bytes.NewReader(runesToUTF8(chars)))

	// skip first char
	r.ReadRune()
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			break
		}

		// detect escape
		if char == '\\' {
			char, _, _ := r.ReadRune()

			// unicode
			if char == 'u' {
				hexs := make([]rune, 4)
				for i := 0; i < 4; i++ {
					char, _, _ := r.ReadRune()
					hexs[i] = char
				}

				dec, _ := strconv.ParseInt(string(hexs), 16, 64)
				newRunes = append(newRunes, rune(dec))
				continue
			}

			// hex
			if char == 'x' {
				hexs := make([]rune, 2)
				for i := 0; i < 2; i++ {
					char, _, _ := r.ReadRune()
					hexs[i] = char
				}

				dec, _ := strconv.ParseInt(string(hexs), 16, 64)
				newRunes = append(newRunes, rune(dec))
				continue
			}

			// line terminator (return carriage or return carriage and line feed)
			if char == '\r' {
				char, _, _ := r.ReadRune()
				if char == '\n' {
					continue
				}

				continue
			}

			// line terminator (line feed)
			if char == '\n' {
				continue
			}

			// line terminator (line separator)
			if char == '\u2028' {
				continue
			}

			// line terminator (paragraph separator)
			if char == '\u2029' {
				continue
			}

			// null
			if char == '0' {
				continue
			}

			// backspace
			if char == 'b' {
				newRunes = append(newRunes, '\b')
				continue
			}

			// form feed
			if char == 'f' {
				newRunes = append(newRunes, '\f')
				continue
			}

			// line feed
			if char == 'n' {
				newRunes = append(newRunes, '\n')
				continue
			}

			// carriage return
			if char == 'r' {
				newRunes = append(newRunes, '\r')
				continue
			}

			// horizontal tab
			if char == 't' {
				newRunes = append(newRunes, '\v')
				continue
			}

			// vertical tab
			if char == 'v' {
				newRunes = append(newRunes, '\v')
				continue
			}

			newRunes = append(newRunes, char)
			continue
		}

		if char == '"' && ty == doubleQuotedStr {
			break
		}

		if char == '\'' && ty == singleQuotedStr {
			break
		}

		newRunes = append(newRunes, char)
	}

	return string(newRunes)
}
