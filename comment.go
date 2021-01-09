package json5extract

import "io"

// comment types
const (
	commInline = iota
	commMultiLine
)

func parseComment(r reader) (int, error) {
	char, _, err := r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 0, ErrInvalidFormat
		}

		return 0, err
	}

	// single line comment
	if char == '/' {
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return 0, ErrInvalidFormat
				}

				return 0, err
			}

			if char == '\r' || char == '\n' {
				break
			}
		}

		return commInline, nil
	}

	// multi line comment
	if char == '*' {
		for {
			char, _, err := r.ReadRune()
			if err != nil {
				if err == io.EOF {
					return 0, ErrInvalidFormat
				}

				return 0, err
			}

			if char == '*' {
				char, _, err := r.ReadRune()
				if err != nil {
					if err == io.EOF {
						return 0, ErrInvalidFormat
					}

					return 0, err
				}

				if char == '/' {
					break
				}
			}
		}

		return commMultiLine, nil
	}

	return 0, ErrInvalidFormat
}
