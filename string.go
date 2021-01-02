package json5extract

import (
	"bytes"
	"strconv"
)

var (
	quote          = byte(39)
	doubleQuote    = byte(34)
	reverseSolidus = byte(92)
	solidus        = byte(47)
	zero           = byte(48)
)

var (
	hex = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F'}
)

// String types
const (
	DoubleQuotedStr = iota
	SingleQuotedStr
	// UnquotedStr only can be used as object key, not value
	// UnquotedStr string cannot contain any escaped character except unicode (\u{hex}{hex}{hex}{hex}).
	// UnquotedStr only valid if first character is latin, followed by any latin character, numeric, unicode, or underscore
	UnquotedStr
)

// String represent JSON5 string type.
type String struct {
	// String type
	Type  int
	Raw   []byte
	Value string
}

// Kind return JSON5 kind
func (str *String) Kind() int {
	return KindStr
}

func (str *String) pushByt(char byte) {
	str.Raw = append(str.Raw, char)
}

func parseStr(r reader, ty int) (*String, error) {
	str := &String{Type: DoubleQuotedStr}

	if ty == DoubleQuotedStr {
		str.pushByt('"')
	} else {
		str.pushByt('\'')
	}

	for {
		char, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		if char == reverseSolidus {
			str.pushByt(char)
			char, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			str.pushByt(char)
			// if unicode (\u), check the hexs values
			if char == 117 {
				for i := 0; i < 4; i++ {
					char, err := r.ReadByte()
					if err != nil {
						return nil, err
					}

					if !isCharHex(char) {
						return nil, ErrInvalidFormat
					}

					str.pushByt(char)
				}

				continue
			}

			// if hexa (\x), check the hexs values
			if char == 120 {
				for i := 0; i < 2; i++ {
					char, err := r.ReadByte()
					if err != nil {
						return nil, err
					}

					if !isCharHex(char) {
						return nil, ErrInvalidFormat
					}

					str.pushByt(char)
				}

				continue
			}

			continue
		}

		if char == '"' && ty == DoubleQuotedStr {
			str.pushByt(char)
			break
		}

		if char == '\'' && ty == SingleQuotedStr {
			str.pushByt(char)
			break
		}

		// line terminator (\n)
		if char == '\n' {
			rawLen := len(str.Raw)
			prev1Char := str.Raw[rawLen-1]
			prev2Char := str.Raw[rawLen-2]
			prev3Char := str.Raw[rawLen-3]

			if prev1Char == reverseSolidus {
				if prev2Char == reverseSolidus {
					return nil, ErrInvalidFormat
				}

				continue
			}

			if prev1Char == '\r' {
				if prev2Char != reverseSolidus {
					return nil, ErrInvalidFormat
				}

				if prev3Char == reverseSolidus {
					return nil, ErrInvalidFormat
				}
			}
		}

		// line terminator (\r)
		if char == '\r' {
			rawLen := len(str.Raw)
			prev1Char := str.Raw[rawLen-1]
			prev2Char := str.Raw[rawLen-2]
			if prev1Char != reverseSolidus {
				return nil, ErrInvalidFormat
			}

			if prev2Char == reverseSolidus {
				return nil, ErrInvalidFormat
			}
		}

		// line terminator (line separator or paragraph separator)
		if char == 226 {
			lineSp := make([]byte, 3)
			lineSp[0] = char
			isLineSp := false
			for i := 0; i < 2; i++ {
				char, _ := r.ReadByte()
				lineSp[i+1] = char

				if i == 0 && char != 128 {
					break
				}

				if i == 1 && char == 168 || char == 169 {
					isLineSp = true
				}
			}

			for _, c := range lineSp {
				str.pushByt(c)
			}

			if isLineSp {
				rawLen := len(str.Raw)
				prev1Char := str.Raw[rawLen-1]
				prev2Char := str.Raw[rawLen-2]
				if prev1Char != reverseSolidus {
					return nil, ErrInvalidFormat
				}

				if prev2Char == reverseSolidus {
					return nil, ErrInvalidFormat
				}

				continue
			}

			continue
		}

		str.pushByt(char)
	}

	str.Value = string(bytes.Trim(unescapeByts(str.Raw), `"`))

	return str, nil
}

func parseSingleQuoteStr(r reader) (*String, error) {
	str := &String{Type: SingleQuotedStr}

	return str, nil
}

func isCharHex(char byte) bool {
	for _, c := range hex {
		if char == c {
			return true
		}
	}

	return false
}

func unescapeByts(byts []byte) []byte {
	newByts := make([]byte, 0)
	r := readFromBytes(byts)
	for {
		char, err := r.ReadByte()
		if err != nil {
			break
		}

		if char == reverseSolidus {
			char, _ := r.ReadByte()

			// unescape unicode
			if char == 'u' {
				hexByts := make([]byte, 0)
				for i := 0; i < 4; i++ {
					char, _ := r.ReadByte()
					hexByts = append(hexByts, char)
				}

				uq, _ := strconv.Unquote(`'\u` + string(hexByts) + `'`)
				newByts = append(newByts, []byte(uq)...)
				continue
			}

			// unescape hex
			if char == 'x' {
				hexByts := make([]byte, 0)
				for i := 0; i < 2; i++ {
					char, _ := r.ReadByte()
					hexByts = append(hexByts, char)
				}

				decimal, _ := strconv.ParseInt(string(hexByts), 16, 64)
				newByts = append(newByts, byte(decimal))
				continue
			}

			// backspace
			if char == 'b' {
				newByts = append(newByts, '\b')
				continue
			}

			// form feed
			if char == 'f' {
				newByts = append(newByts, '\f')
				continue
			}

			// line feed
			if char == 'n' {
				newByts = append(newByts, '\n')
				continue
			}

			// carriage return
			if char == 'r' {
				newByts = append(newByts, '\r')
				continue
			}

			// horizontal tab
			if char == 't' {
				newByts = append(newByts, '\t')
				continue
			}

			// vertical tab
			if char == 'v' {
				newByts = append(newByts, '\v')
				continue
			}

			// null
			if char == '0' {
				newByts = append(newByts, []byte(`\0`)...)
				continue
			}

			// line terminator (\r with \n or just \r)
			if char == '\r' {
				char, _ := r.ReadByte()
				if char != '\n' {
					newByts = append(newByts, char)
				}

				continue
			}

			// line terminator (\n)
			if char == '\n' {
				continue
			}

			// line terminator (line separator or paragraph separator)
			if char == 226 {
				lineSp := make([]byte, 3)
				lineSp[0] = char
				isLineSp := false
				for i := 0; i < 2; i++ {
					char, _ := r.ReadByte()
					lineSp[i+1] = char

					if i == 0 && char != 128 {
						break
					}

					if i == 1 && char == 168 || char == 169 {
						isLineSp = true
					}
				}

				if isLineSp {
					continue
				} else {
					for _, c := range lineSp {
						newByts = append(newByts, c)
					}
				}

				continue
			}

			newByts = append(newByts, char)
			continue
		}

		newByts = append(newByts, char)
	}

	return newByts
}
