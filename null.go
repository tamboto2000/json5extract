package json5extract

// continuation after 'n' suspected as null
var null = []rune("ull")

func parseNull(r reader) (*JSON5, error) {
	nll := &JSON5{kind: Null}
	nll.push('n')
	for _, c := range null {
		char, _, err := r.ReadRune()
		if err != nil {
			return nil, ErrInvalidFormat
		}

		if char != c {
			return nil, ErrInvalidFormat
		}

		nll.push(char)
	}

	return nll, nil
}
