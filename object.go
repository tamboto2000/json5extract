package json5extract

import (
	"io"
	"unicode"
)

func parseObj(r reader) (*JSON5, error) {
	obj := &JSON5{Kind: Object, val: make(map[string]*JSON5)}
	obj.push('{')
	isEnd, err := parseKeyVal(r, obj)
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

			obj.push(char)
			onNext = true
			continue
		}

		if char == '}' {
			obj.push(char)
			break
		}

		// comment
		if char == '/' {
			if err := parseComment(r); err != nil {
				return nil, err
			}

			continue
		}

		r.UnreadRune()
		isEnd, err := parseKeyVal(r, obj)
		if err != nil {
			return nil, err
		}

		if isEnd {
			break
		}

		onNext = false
	}

	return obj, nil
}

func parseKeyVal(r reader, obj *JSON5) (isEnd bool, err error) {
	keyVal := obj.val.(map[string]*JSON5)
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return isEnd, ErrInvalidFormat
			}

			return isEnd, err
		}

		if char == '}' {
			obj.push(char)
			return true, nil
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char == '/' {
			if err := parseComment(r); err != nil {
				return isEnd, err
			}

			continue
		}

		r.UnreadRune()
		if id, str, err := parseIdentifier(r); err == nil {
			obj.pushRns([]rune(str + ":"))
			val, err := parse(r)
			if err != nil {
				return isEnd, err
			}

			keyVal[id] = val
			obj.pushRns(val.raw)
			obj.val = keyVal
			break
		} else {
			return isEnd, err
		}
	}

	return isEnd, nil
}

func parseComment(r reader) error {
	char, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return ErrInvalidFormat
		}

		return err
	}

	// single line comment
	if char == '/' {
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return ErrInvalidFormat
				}

				return err
			}

			if char == '\r' || char == '\n' {
				break
			}
		}

		return nil
	}

	// multi line comment
	if char == '*' {
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return ErrInvalidFormat
				}

				return err
			}

			if char == '*' {
				char, _, err := r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return ErrInvalidFormat
					}

					return err
				}

				if char != '/' {
					return ErrInvalidFormat
				}

				break
			}
		}

		return nil
	}

	return ErrInvalidFormat
}
