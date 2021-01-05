package json5extract

import (
	"io"
	"math"
	"strconv"
	"unicode"
)

// Number types
const (
	Int = iota
	Float
	Infinity
	NaN
)

// Number represent JSON5 number
type Number struct {
	Type     int
	IntValue int64
	// Infinity and NaN is stored in FloatValue field.
	FloatValue float64
	// Hexadecimal number (0{x or X}{hex1}{hex2}...{hexn})
	IsHex bool
	//Exponent number ({num1}{num2}...{numn}{e or E}{num1}{num2}...{numn})
	WithExp    bool
	IsPositive bool
	IsNegative bool
	raw        []rune
}

var (
	// continuation after 'I' suspected as Infinity
	inf = []rune("nfinity")
	// continuation after 'N' suspected as NaN
	nan = []rune("aN")
)

// Kind return JSON5 kind
func (num *Number) Kind() int {
	return KindNum
}

// RawBytes return raw JSON5 to []byte
func (num *Number) RawBytes() []byte {
	return runesToUTF8(num.raw)
}

// RawRunes return raw JSON5 to []rune
func (num *Number) RawRunes() []rune {
	return num.raw
}

func (num *Number) push(char rune) {
	num.raw = append(num.raw, char)
}

func parseNum(r reader, firstC rune) (*Number, error) {
	num := &Number{Type: Int, IsPositive: true}
	num.push(firstC)

	if isMinOrPlusSign(firstC) {
		if firstC == '-' {
			num.IsPositive = false
			num.IsNegative = true
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		num.push(char)
		firstC = char
	}

	if !isCharNumBegin(firstC) {
		return nil, ErrInvalidFormat
	}

	if firstC == '0' {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return num, nil
			}

			return nil, err
		}

		if char == ',' || char == '}' || char == ']' {
			return num, nil
		}

		if char == '.' {
			num.push(char)
			num.Type = Float

			char, _, err := r.ReadRune()
			if err != nil {
				return nil, err
			}

			num.push(char)

			if !unicode.IsNumber(char) {
				return nil, ErrInvalidFormat
			}

			if err := parseOnlyNum(r, num); err != nil {
				return nil, err
			}

			return num, nil
		}

		if char == 'e' || char == 'E' {
			num.push(char)
			if err := parseExp(r, num); err != nil {
				return nil, err
			}

			if err := parseOnlyNum(r, num); err != nil {
				return nil, err
			}

			return num, nil
		}

		// detect hexa number
		if char == 'x' || char == 'X' {
			num.push(char)
			char, _, err := r.ReadRune()
			if err != nil {
				return nil, err
			}

			num.push(char)
			if !isCharHex(char) {
				return nil, ErrInvalidFormat
			}

			num.IsHex = true
			if err := parseOnlyHex(r, num); err != nil {
				return nil, err
			}

			return num, nil
		}

		if char == ']' || char == '}' || char == ',' {
			r.UnreadRune()
			return num, nil
		}

		if unicode.IsControl(char) || char == ' ' {
			return num, nil
		}

		return nil, ErrInvalidFormat
	}

	if firstC == '.' {
		char, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		if !unicode.IsNumber(char) {
			return nil, ErrInvalidFormat
		}

		num.push(char)
		num.Type = Float

		if err := parseOnlyNum(r, num); err != nil {
			return nil, err
		}
	}

	if unicode.IsNumber(firstC) {
		if err := parseOnlyNum(r, num); err != nil {
			return nil, err
		}
	}

	if firstC == 'I' {
		if err := parseInf(r, num); err != nil {
			return nil, err
		}

		num.Type = Infinity
	}

	if firstC == 'N' {
		if err := parseNaN(r, num); err != nil {
			return nil, err
		}

		num.Type = NaN
	}

	return num, nil
}

func parseOnlyNum(r reader, num *Number) error {
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if !unicode.IsNumber(char) {
			if char == '.' {
				if num.Type == Float {
					return ErrInvalidFormat
				}

				num.push(char)
				num.Type = Float
				continue
			}

			if char == 'e' || char == 'E' {
				if num.WithExp {
					return ErrInvalidFormat
				}

				num.push(char)

				if err := parseExp(r, num); err != nil {
					return err
				}

				num.WithExp = true
				continue
			}

			if char == ']' || char == '}' || char == ',' {
				r.UnreadRune()
				break
			}

			if unicode.IsControl(char) || char == ' ' {
				break
			}

			return ErrInvalidFormat
		}

		num.push(char)
	}

	if num.Type == Float {
		i, _ := strconv.ParseFloat(string(num.raw), 64)
		num.FloatValue = i
	}

	if num.Type == Int {
		if num.WithExp {
			num.Type = Float
			i, _ := strconv.ParseFloat(string(num.raw), 64)
			num.FloatValue = i
		}

		i, _ := strconv.ParseInt(string(num.raw), 10, 64)
		num.IntValue = i
	}

	return nil
}

func parseOnlyHex(r reader, num *Number) error {
	for {
		char, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if !isCharHex(char) {
			if char == ']' || char == '}' || char == ',' {
				r.UnreadRune()
				break
			}

			if unicode.IsControl(char) || char == ' ' {
				break
			}

			return ErrInvalidFormat
		}

		num.push(char)
	}

	i, _ := strconv.ParseInt(string(num.raw), 0, 64)
	num.IntValue = i

	return nil
}

// parse exponent
func parseExp(r reader, num *Number) error {
	for i := 0; i < 2; i++ {
		char, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		num.push(char)
		if i == 0 {
			if isMinOrPlusSign(char) {
				continue
			}

			if unicode.IsNumber(char) {
				break
			}

			if isCharHex(char) {
				if !num.IsHex {
					return ErrInvalidFormat
				}

				break
			}

			return ErrInvalidFormat
		}

		if i == 1 {
			if unicode.IsNumber(char) {
				break
			}

			if isCharHex(char) {
				if !num.IsHex {
					return ErrInvalidFormat
				}

				break
			}

			return ErrInvalidFormat
		}
	}

	num.WithExp = true
	return nil
}

func parseInf(r reader, num *Number) error {
	for _, c := range inf {
		char, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		if char != c {
			return ErrInvalidFormat
		}

		num.push(char)
	}

	char := num.raw[0]
	if char == '-' {
		num.FloatValue = math.Inf(-0)
	} else {
		num.FloatValue = math.Inf(0)
	}

	return nil
}

func parseNaN(r reader, num *Number) error {
	for _, c := range nan {
		char, _, err := r.ReadRune()
		if err != nil {
			return err
		}

		if char != c {
			return ErrInvalidFormat
		}

		num.push(char)
	}

	num.FloatValue = math.NaN()

	return nil
}

func isMinOrPlusSign(char rune) bool {
	if char == '-' || char == '+' {
		return true
	}

	return false
}
