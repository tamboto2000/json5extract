package json5extract

import (
	"io"
	"unicode"
)

func parseArray(r reader) (*JSON5, error) {
	arr := &JSON5{Kind: Array}
	arr.push('[')
	vals := make([]*JSON5, 0)

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

		if char == ']' {
			arr.push(char)
			return arr, nil
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return nil, ErrInvalidFormat
			}

			continue
		}

		json5, err := parse(r, char)
		if err != nil {
			return nil, err
		}

		if json5 != nil {
			vals = append(vals, json5)
			arr.pushRns(json5.RawRunes())
			break
		}

		return nil, ErrInvalidFormat
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

			arr.push(',')
			onNext = true
			continue
		}

		if char == ']' {
			arr.push(char)
			break
		}

		// comment
		if char == '/' {
			if _, err := parseComment(r); err != nil {
				return nil, ErrInvalidFormat
			}

			continue
		}

		if onNext {
			json, err := parse(r, char)
			if err != nil {
				return nil, err
			}

			if json != nil {
				arr.pushRns(json.raw)
				vals = append(vals, json)
				onNext = false
				continue
			}

			return nil, ErrInvalidFormat
		}

		return nil, ErrInvalidFormat
	}

	arr.val = vals
	return arr, nil
}
