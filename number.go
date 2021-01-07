package json5extract

import (
	"io"
	"math"
	"strconv"
	"unicode"
)

type numStates struct {
	isInt      bool
	isFloat    bool
	isHex      bool
	withExp    bool
	isPositive bool
	isNegative bool
	isInfinity bool
	isNan      bool
}

var (
	// continuation after 'I' suspected as Infinity
	inf = []rune("nfinity")
	// continuation after 'N' suspected as NaN
	nan = []rune("aN")
)

func parseNum(r reader, firstC rune) (*JSON5, error) {
	num := new(JSON5)
	num.push(firstC)
	state := new(numStates)

	state.isPositive = true
	state.isInt = true

	if isMinOrPlusSign(firstC) {
		if firstC == '-' {
			state.isPositive = false
			state.isNegative = true
		}

		char, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		if isMinOrPlusSign(char) {
			return nil, ErrInvalidFormat
		}

		if unicode.IsControl(char) || char == ' ' {
			return nil, ErrInvalidFormat
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
			r.UnreadRune()
			return num, nil
		}

		if char == '.' {
			num.push(char)
			state.isFloat = true
			state.isInt = false

			char, _, err := r.ReadRune()
			if err != nil {
				return nil, err
			}

			num.push(char)

			if !unicode.IsNumber(char) {
				return nil, ErrInvalidFormat
			}

			if err := parseOnlyNum(r, num, state); err != nil {
				return nil, err
			}

			return num, nil
		}

		if char == 'e' || char == 'E' {
			num.push(char)
			if err := parseExp(r, num, state); err != nil {
				return nil, err
			}

			if err := parseOnlyNum(r, num, state); err != nil {
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

			state.isHex = true
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
		state.isFloat = true
		state.isInt = false

		if err := parseOnlyNum(r, num, state); err != nil {
			return nil, err
		}
	}

	if unicode.IsNumber(firstC) {
		if err := parseOnlyNum(r, num, state); err != nil {
			return nil, err
		}
	}

	if firstC == 'I' {
		if err := parseInf(r, num, state); err != nil {
			return nil, err
		}

		state.isInfinity = true
	}

	if firstC == 'N' {
		if err := parseNaN(r, num, state); err != nil {
			return nil, err
		}

		state.isNan = true
	}

	return num, nil
}

func parseOnlyNum(r reader, num *JSON5, state *numStates) error {
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
				if state.isFloat {
					return ErrInvalidFormat
				}

				num.push(char)
				state.isFloat = true
				state.isInt = false
				continue
			}

			if char == 'e' || char == 'E' {
				if state.withExp {
					return ErrInvalidFormat
				}

				num.push(char)

				if err := parseExp(r, num, state); err != nil {
					return err
				}

				state.withExp = true
				state.isFloat = true
				state.isInt = false
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

	if state.isFloat {
		i, _ := strconv.ParseFloat(string(num.raw), 64)
		num.Kind = Float
		num.val = i
	}

	if state.isInt {
		if state.withExp {
			i, _ := strconv.ParseFloat(string(num.raw), 64)
			num.Kind = Float
			num.val = i
		}

		i, _ := strconv.ParseInt(string(num.raw), 10, 64)
		num.Kind = Integer
		num.val = i
	}

	return nil
}

func parseOnlyHex(r reader, num *JSON5) error {
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
	num.val = i
	num.Kind = Integer

	return nil
}

// parse exponent
func parseExp(r reader, num *JSON5, state *numStates) error {
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
				if !state.isHex {
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
				if !state.isHex {
					return ErrInvalidFormat
				}

				break
			}

			return ErrInvalidFormat
		}
	}

	state.withExp = true
	return nil
}

func parseInf(r reader, num *JSON5, state *numStates) error {
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

	if state.isNegative {
		num.val = math.Inf(-1)
	} else {
		num.val = math.Inf(1)
	}

	num.Kind = Infinity

	return nil
}

func parseNaN(r reader, num *JSON5, state *numStates) error {
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

	if state.isNegative {
		num.val = -math.NaN()
	} else {
		num.val = math.NaN()
	}

	num.Kind = NaN

	return nil
}

func isMinOrPlusSign(char rune) bool {
	if char == '-' || char == '+' {
		return true
	}

	return false
}
