package json5extract

import (
	"io"
	"unicode"
)

// This file contains parser method for unquoted string object identifier. This string can't be used as value, only
// can be used as object identifier

// This is list of reserved words that can't be used as unquoted object identifier name,
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

// Identifier name obey the  ECMAScript 5.1 Lexical Grammar, see
// https://www.ecma-international.org/ecma-262/5.1/#sec-7.6 "Identifier Names and Identifiers"
func parseIdentifier(r reader, char rune) (id, raw []rune, err error) {
	rs := make([]rune, 0)

	// double quoted string
	if char == '"' {
		// find key
		str, err := parseStr(r, doubleQuotedStr)
		if err != nil {
			return nil, nil, err
		}

		// find key terminator
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, nil, ErrInvalidFormat
				}

				return nil, nil, err
			}

			if unicode.IsControl(char) || char == ' ' {
				continue
			}

			if char == '/' {
				if _, err := parseComment(r); err != nil {
					return nil, nil, ErrInvalidFormat
				}

				continue
			}

			if char != ':' {
				return nil, nil, ErrInvalidFormat
			}

			break
		}

		return []rune(str.val.(string)), str.raw, nil
	}

	// single quoted string
	if char == '\'' {
		// find key
		str, err := parseStr(r, singleQuotedStr)
		if err != nil {
			return nil, nil, err
		}

		// find key terminator
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, nil, ErrInvalidFormat
				}

				return nil, nil, err
			}

			if unicode.IsControl(char) || char == ' ' {
				continue
			}

			// comment
			if char == '/' {
				if _, err := parseComment(r); err != nil {
					return nil, nil, ErrInvalidFormat
				}

				continue
			}

			if char != ':' {
				return nil, nil, ErrInvalidFormat
			}

			break
		}

		return []rune(str.val.(string)), str.raw, nil
	}

	// unicode
	if char == '\\' {
		rn, err := parseUnicode(r)
		if err != nil {
			return nil, nil, err
		}

		char = rn
	}

	if isCharIDValid(char, true) {
		rs = append(rs, char)
		// extract key name
		isIDEnd := false
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return nil, nil, ErrInvalidFormat
				}

				return nil, nil, err
			}

			// find terminator
			if isIDEnd {
				if unicode.IsControl(char) || char == ' ' {
					isIDEnd = true
					continue
				}

				if char != ':' {
					return nil, nil, ErrInvalidFormat
				}

				break
			}

			// unicode
			if char == '\\' {
				rn, err := parseUnicode(r)
				if err != nil {
					return nil, nil, err
				}

				rs = append(rs, rn)
				continue
			}

			// comment
			if char == '/' {
				if _, err := parseComment(r); err != nil {
					return nil, nil, ErrInvalidFormat
				}

				continue
			}

			if unicode.IsControl(char) || char == ' ' {
				isIDEnd = true
				continue
			}

			// key terminator
			if char == ':' {
				break
			}

			if isCharIDValid(char, false) {
				rs = append(rs, char)
				continue
			}

			return nil, nil, ErrInvalidFormat
		}

		if isReservedWord(rs) {
			return nil, nil, ErrInvalidFormat
		}

		return rs, rs, nil
	}

	return nil, nil, ErrInvalidFormat
}

func isCharIDValid(char rune, begin bool) bool {
	if unicode.IsLetter(char) {
		return true
	}

	if unicode.IsNumber(char) {
		if begin {
			return false
		}

		return true
	}

	if char == '$' {
		return true
	}

	if char == '_' {
		return true
	}

	return false
}
