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

		r.UnreadRune()
		json5, err := parse(r)
		if err != nil {
			return nil, err
		}

		vals = append(vals, json5)
		arr.pushRns(json5.RawRunes())
		break
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

			arr.push(char)
			onNext = true
			continue
		}

		if char == ']' {
			arr.push(char)
			break
		}

		r.UnreadRune()
		json5, err := parse(r)
		if err != nil {
			return nil, err
		}

		vals = append(vals, json5)
		arr.pushRns(json5.RawRunes())
		onNext = false
	}

	arr.val = vals
	return arr, nil
}
