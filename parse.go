package json5extract

import (
	"errors"
	"io"
)

// JSON5 kinds
const (
	String = iota
	Integer
	Float
	Infinity
	NaN
	Boolean
	Null
	Array
	Object
)

// JSON5 represent parsed value of JSON5 types. Check JSON5.Kind to know
// which data type a value is
type JSON5 struct {
	Kind int
	val  interface{}
	raw  []rune
}

// String return string value
func (json *JSON5) String() (string, error) {
	if json.Kind != String {
		return "", errors.New("value is not string")
	}

	return json.val.(string), nil
}

// Integer return int64 value
func (json *JSON5) Integer() (int64, error) {
	if json.Kind != Integer {
		return 0, errors.New("value is not int")
	}

	return json.val.(int64), nil
}

// Float return float64 value
func (json *JSON5) Float() (float64, error) {
	if json.Kind != Float {
		return 0, errors.New("value is not float")
	}

	return json.val.(float64), nil
}

// Infinity return Inf value
func (json *JSON5) Infinity() (float64, error) {
	if json.Kind != Infinity {
		return 0, errors.New("value is not Inf")
	}

	return json.val.(float64), nil
}

// NaN return NaN value
func (json *JSON5) NaN() (float64, error) {
	if json.Kind != NaN {
		return 0, errors.New("value is not NaN")
	}

	return json.val.(float64), nil
}

// Boolean return bool value
func (json *JSON5) Boolean() (bool, error) {
	if json.Kind != Boolean {
		return false, errors.New("value is not bool")
	}

	return json.val.(bool), nil
}

// Null will not return any value because, well, is is nil...
func (json *JSON5) Null() error {
	if json.Kind != Null {
		return errors.New("value is not null")
	}

	return nil
}

// Array return slice of JSON5 values
func (json *JSON5) Array() ([]*JSON5, error) {
	if json.Kind != Array {
		return nil, errors.New("value is not array")
	}

	return json.val.([]*JSON5), nil
}

// Object return map of JSON5 values
func (json *JSON5) Object() (map[string]*JSON5, error) {
	if json.Kind != Object {
		return nil, errors.New("value is not object")
	}

	return json.val.(map[string]*JSON5), nil
}

// RawBytes return parsed raw bytes of JSON5
func (json *JSON5) RawBytes() []byte {
	return runesToUTF8(json.raw)
}

// RawRunes return parsed raw bytes of JSON5
func (json *JSON5) RawRunes() []rune {
	return json.raw
}

func (json *JSON5) push(char rune) {
	json.raw = append(json.raw, char)
}

func (json *JSON5) pushRns(chars []rune) {
	json.raw = append(json.raw, chars...)
}

func parseAll(r reader) ([]*JSON5, error) {
	json5s := make([]*JSON5, 0)
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

func parse(r reader) (*JSON5, error) {
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
			json5, err := parseStr(r, doubleQuotedStr)
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
			json5, err := parseStr(r, singleQuotedStr)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse number
		if isCharNumBegin(char) {
			json5, err := parseNum(r, char)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse true boolean
		if char == 't' {
			json5, err := parseTrueBool(r)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse false boolean
		if char == 'f' {
			json5, err := parseFalseBool(r)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse null
		if char == 'n' {
			json5, err := parseNull(r)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}

			return json5, nil
		}

		// parse array
		if char == '[' {
			json5, err := parseArray(r)
			if err != nil {
				if err == io.EOF {
					break
				}

				r.UnreadRune()
				return nil, err
			}
			return json5, nil
		}

		// parse object
		if char == '{' {
			json5, err := parseObj(r)
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
