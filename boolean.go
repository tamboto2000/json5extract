package json5extract

var (
	// continuation after 'f' suspected as false boolean
	boolFalse = []rune("alse")
	// continuation after 't' suspected as true boolean
	boolTrue = []rune("rue")
)

// Boolean represent JSON5 boolean
type Boolean struct {
	raw   []rune
	Value bool
}

func (bl *Boolean) Kind() int {
	return KindBool
}

func (bl *Boolean) RawRunes() []rune {
	return bl.raw
}

func (bl *Boolean) RawBytes() []byte {
	return runesToUTF8(bl.raw)
}

func (bl *Boolean) push(char rune) {
	bl.raw = append(bl.raw, char)
}

func parseTrueBool(r reader) (*Boolean, error) {
	bl := new(Boolean)
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

	bl.Value = true

	return bl, nil
}

func parseFalseBool(r reader) (*Boolean, error) {
	bl := new(Boolean)
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

	bl.Value = false

	return bl, nil
}
