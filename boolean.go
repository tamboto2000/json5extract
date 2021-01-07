package json5extract

var (
	// continuation after 'f' suspected as false boolean
	boolFalse = []rune("alse")
	// continuation after 't' suspected as true boolean
	boolTrue = []rune("rue")
)

func parseTrueBool(r reader) (*JSON5, error) {
	bl := &JSON5{Kind: Boolean}
	bl.push('t')
	for _, c := range boolTrue {
		char, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		if char != c {
			return nil, ErrInvalidFormat
		}

		bl.push(char)
	}

	bl.val = true

	return bl, nil
}

func parseFalseBool(r reader) (*JSON5, error) {
	bl := &JSON5{Kind: Boolean}
	bl.push('f')
	for _, c := range boolFalse {
		char, _, err := r.ReadRune()
		if err != nil {
			return nil, ErrInvalidFormat
		}

		if char != c {
			return nil, ErrInvalidFormat
		}

		bl.push(char)
	}

	bl.val = false

	return bl, nil
}
