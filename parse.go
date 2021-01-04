package json5extract

import "io"

// JSON5 kinds
const (
	KindStr = iota
)

// JSON5 is interface for parsed JSON5 type
type JSON5 interface {
	// Kind return JSON5 kind
	Kind() int
	// Return raw JSON5 in []byte format
	RawBytes() []byte
	// Return raw JSON5 in []rune format
	RawRunes() []rune
}

func parseAll(r reader) ([]JSON5, error) {
	json5s := make([]JSON5, 0)
	for {
		json5, err := parse(r)
		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		json5s = append(json5s, json5)
	}

	return json5s, nil
}

func parse(r reader) (JSON5, error) {
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		// parse double quoted string
		if char == '"' {
			json5, err := parseStr(r, DoubleQuotedStr)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse single quoted string
		if char == '\'' {
			json5, err := parseStr(r, SingleQuotedStr)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}
	}

	return nil, io.EOF
}
