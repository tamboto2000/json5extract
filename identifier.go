package json5extract

import (
	"strconv"
	"unicode"
)

// import (
// 	"strconv"
// 	"unicode"
// )

// This file contains parser method for unquoted string object identifier. This string can't be used as value, only
// can be used as object identifier

// This is list of reserved words that can't be used as object identifier name,
// also include future reserved words, and a list of considered future reserved words
// References:
// Reserved words
// https://www.ecma-international.org/ecma-262/5.1/#sec-7.6.1.1
// Future reserved words
// https://www.ecma-international.org/ecma-262/5.1/#sec-7.6.1.2
var rsvWords = [][]rune{
	[]rune("break"), []rune("case"),
	[]rune("catch"), []rune("continue"),
	[]rune("debugger"), []rune("default"),
	[]rune("delete"), []rune("do"),
	[]rune("else"), []rune("finally"),
	[]rune("for"), []rune("function"),
	[]rune("if"), []rune("in"),
	[]rune("instanceof"), []rune("new"),
	[]rune("return"), []rune("switch"),
	[]rune("this"), []rune("throw"),
	[]rune("try"), []rune("typeof"),
	[]rune("var"), []rune("void"),
	[]rune("while"), []rune("with"),
	[]rune("class"), []rune("const"),
	[]rune("enum"), []rune("export"),
	[]rune("extends"), []rune("import"),
	[]rune("super"),
}

var ctrlChars = []rune{'\b', '\f', '\n', '\r', '\t', '\v', ' '}

// TestParseIdentifier test identifier name validity and return identifier name
func TestParseIdentifier(byts []byte) (string, bool) {
	r := readFromBytes(byts)
	id, err := parseIdentifier(r)
	if err != nil {
		return "", false
	}

	return id, true
}

// Identifier name obey the  ECMAScript 5.1 Lexical Grammar, see
// https://www.ecma-international.org/ecma-262/5.1/#sec-7.6 "Identifier Names and Identifiers"
func parseIdentifier(r reader) (string, error) {
	rs := make([]rune, 0)
	isIDEnd := false

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		if unicode.IsControl(char) {
			continue
		}

		if char == '"' {
			str, err := parseStr(r, DoubleQuotedStr)
			if err != nil {
				return "", err
			}

			// find terminator
			for {
				char, _, err := r.ReadRune()
				if err != nil {
					return "", err
				}

				if unicode.IsControl(char) || char == ' ' {
					continue
				}

				if char != ':' {
					return "", ErrInvalidFormat
				}

				break
			}

			return str.Value, nil
		}

		if char == '\'' {
			str, err := parseStr(r, SingleQuotedStr)
			if err != nil {
				return "", err
			}

			// find terminator
			for {
				char, _, err := r.ReadRune()
				if err != nil {
					return "", err
				}

				if unicode.IsControl(char) || char == ' ' {
					continue
				}

				if char != ':' {
					return "", ErrInvalidFormat
				}

				break
			}

			return str.Value, nil
		}

		r.UnreadRune()
		break
	}

	for {
		char, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}

		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			if char == '$' {
				rs = append(rs, char)
				continue
			}

			if char == '_' {
				rs = append(rs, char)
				continue
			}

			if char == '\\' {
				char, _, err := r.ReadRune()
				if err != nil {
					return "", err
				}

				if char != 'u' {
					return "", ErrInvalidFormat
				}

				hexs := make([]rune, 4)
				for i := 0; i < 4; i++ {
					char, _, err := r.ReadRune()
					if err != nil {
						return "", err
					}

					hexs[i] = char
				}

				dec, err := strconv.ParseInt(string(hexs), 16, 64)
				if err != nil {
					return "", ErrInvalidFormat
				}

				decRune := rune(dec)
				if !unicode.IsLetter(decRune) && !unicode.IsNumber(decRune) &&
					decRune != '$' && decRune != '_' {
					return "", ErrInvalidFormat
				}

				rs = append(rs, decRune)

				continue
			}

			// identifier terminator
			if char == ':' {
				break
			}

			// identifier end
			if unicode.IsControl(char) || char == ' ' {
				isIDEnd = true
				continue
			}

			if isIDEnd {
				if char != ':' {
					return "", ErrInvalidFormat
				}

				break
			}

			return "", ErrInvalidFormat
		}

		rs = append(rs, char)
	}

	if isReservedWord(rs) {
		return "", ErrInvalidFormat
	}

	return string(rs), nil
}
