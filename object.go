package json5extract

import (
	"io"
	"unicode"
)

func parseObj(r reader) (*JSON5, error) {
	obj := &JSON5{kind: Object, val: make(map[string]*JSON5)}
	state := new(objState)
	obj.push('{')
	err := parseKeyVal(r, obj, state)
	if err != nil {
		return nil, err
	}

	if state.isEnd {
		return obj, nil
	}

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
			if state.onNext {
				return nil, ErrInvalidFormat
			}

			obj.push(',')

			state.onNext = true
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
		if state.onNext {
			err := parseKeyVal(r, obj, state)
			if err != nil {
				return nil, err
			}

			if state.isEnd {
				break
			}

			state.onNext = false
			continue
		}

		return nil, ErrInvalidFormat
	}

	return obj, nil
}

type objState struct {
	isEnd  bool
	onNext bool
}

func parseKeyVal(r reader, obj *JSON5, state *objState) error {
	keyVal := obj.val.(map[string]*JSON5)
	var id []rune
	var idRaw []rune
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return ErrInvalidFormat
			}

			return err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		if char == '}' {
			obj.push(char)
			state.isEnd = true
			return nil
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return err
			}

			continue
		}

		i, iraw, err := parseIdentifier(r, char)
		if err != nil {
			return err
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
				return ErrInvalidFormat
			}

			return err
		}

		if unicode.IsControl(char) || char == ' ' {
			continue
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return err
			}

			continue
		}

		val, err := parse(r, char)
		if err != nil {
			return err
		}

		if val != nil {
			keyVal[string(id)] = val
			obj.pushRns(idRaw)
			obj.pushRns(val.raw)
			obj.val = keyVal

			break
		}

		return ErrInvalidFormat
	}

	return nil
}
