package json5extract

import (
	"io"
	"strconv"
)

func parseUnicode(r reader) (rune, error) {
	char, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0, ErrInvalidFormat
		}

		return 0, err
	}

	if char != 'u' {
		return 0, ErrInvalidFormat
	}

	hexs := make([]rune, 4)
	for i := 0; i < 4; i++ {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return 0, ErrInvalidFormat
			}

			return 0, err
		}

		if !isCharHex(char) {
			return 0, ErrInvalidFormat
		}

		hexs[i] = char
	}

	i, err := strconv.ParseInt(string(hexs), 16, 32)

	return rune(i), err
}
