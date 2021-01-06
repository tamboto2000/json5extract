package json5extract

// continuation after 'n' suspected as null
var null = []rune("ull")

// Null represent JSON5 null value
type Null struct {
	raw []rune
}

func (null *Null) Kind() int {
	return KindNull
}

func (null *Null) RawRunes() []rune {
	return null.raw
}

func (null *Null) RawBytes() []byte {
	return runesToUTF8(null.raw)
}

func (null *Null) push(char rune) {
	null.raw = append(null.raw, char)
}

func parseNull(r reader) (*Null, error) {
	nll := new(Null)
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
