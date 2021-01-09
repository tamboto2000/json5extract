package json5extract

import (
	"io"
	"unicode"
)

func parseObj(r reader) (*JSON5, error) {
	obj := &JSON5{Kind: Object, val: make(map[string]*JSON5)}
	obj.push('{')
	isEnd, err := parseKeyVal(r, obj, true)
	if err != nil {
		return nil, err
	}

	if isEnd {
		return obj, nil
	}

	onNext := false
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil, ErrInvalidFormat
			}

			return nil, err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char == ',' {
			if onNext {
				return nil, ErrInvalidFormat
			}

			onNext = true
			continue
		}

		if char == '}' {
			obj.push(char)
			break
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return nil, err
			}

			continue
		}

		r.UnreadRune()
		if onNext {
			isEnd, err := parseKeyVal(r, obj, false)
			if err != nil {
				return nil, err
			}

			if isEnd {
				break
			}

			onNext = false
			continue
		}

		return nil, ErrInvalidFormat
	}

	return obj, nil
}

func parseKeyVal(r reader, obj *JSON5, isFirst bool) (bool, error) {
	keyVal := obj.val.(map[string]*JSON5)
	var id []rune
	var idRaw []rune
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return false, ErrInvalidFormat
			}

			return false, err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char == '}' {
			if !isFirst {
				obj.push(',')
			}

			obj.push(char)
			return true, nil
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return false, err
			}

			continue
		}

		i, iraw, err := parseIdentifier(r, char)
		if err != nil {
			return false, err
		}

		id = i
		idRaw = iraw
		idRaw = append(idRaw, ':')
		break
	}

	// find value
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return false, ErrInvalidFormat
			}

			return false, err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return false, err
			}

			continue
		}

		val, err := parse(r, char)
		if err != nil {
			return false, err
		}

		if val != nil {
			keyVal[string(id)] = val
			obj.pushRns(idRaw)
			obj.pushRns(val.raw)
			obj.push(',')
			obj.val = keyVal
		}

		break
	}

	return false, nil
}

// comment types
const (
	commInline = iota
	commMultiLine
)
