package json5extract

import (
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
	kind int
	val  interface{}
	raw  []rune
}

// Kind return json kind
func (json *JSON5) Kind() int {
	return json.kind
}

// String return string value. Will panic if kind is not String
func (json *JSON5) String() string {
	if json.kind != String {
		panic("value is not string")
	}

	return json.val.(string)
}

// Integer return int64 value. Will panic if kind is not Integer
func (json *JSON5) Integer() int64 {
	if json.kind != Integer {
		panic("value is not int")
	}

	return json.val.(int64)
}

// Float return float64 value. Will panic if kind is not Float
func (json *JSON5) Float() float64 {
	if json.kind != Float {
		panic("value is not float")
	}

	return json.val.(float64)
}

// Infinity return Inf value. Will panic if kind is not Infinity
func (json *JSON5) Infinity() float64 {
	if json.kind != Infinity {
		panic("value is not Inf")
	}

	return json.val.(float64)
}

// NaN return NaN value. Will panic if kind is not NaN
func (json *JSON5) NaN() float64 {
	if json.kind != NaN {
		panic("value is not NaN")
	}

	return json.val.(float64)
}

// Boolean return bool value. Will panic if kind is not Boolean
func (json *JSON5) Boolean() bool {
	if json.kind != Boolean {
		panic("value is not bool")
	}

	return json.val.(bool)
}

// Array return slice of JSON5 values. Will panic if kind is not Array
func (json *JSON5) Array() []*JSON5 {
	if json.kind != Array {
		panic("value is not array")
	}

	return json.val.([]*JSON5)
}

// Object return map of JSON5 values. Will panic if kind is not Object
func (json *JSON5) Object() map[string]*JSON5 {
	if json.kind != Object {
		panic("value is not object")
	}

	return json.val.(map[string]*JSON5)
}

// Bytes return parsed raw bytes of JSON5
func (json *JSON5) Bytes() []byte {
	return runesToUTF8(json.raw)
}

// Runes return parsed raw bytes of JSON5
func (json *JSON5) Runes() []rune {
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
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		json5, err := parse(r, char)
		if err != nil {
			if err == io.EOF {
				break
			}

			continue
		}

		if json5 != nil {
			json5s = append(json5s, json5)
		}
	}

	return json5s, nil
}

func parse(r reader, char rune) (*JSON5, error) {
	// parse double quoted string
	if char == '"' {
		json5, err := parseStr(r, doubleQuotedStr)
		if err != nil {
			if err == io.EOF {
				return nil, err
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
				return nil, err
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
				return nil, err
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
				return nil, err
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
				return nil, err
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
				return nil, err
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
				return nil, err
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
				return nil, err
			}

			r.UnreadRune()
			return nil, err
		}
		return json5, nil
	}

	return nil, nil
}
