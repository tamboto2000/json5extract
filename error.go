package json5extract

import "errors"

// ErrInvalidFormat occured when a data is invalid format, such as unquoted string with hex escape (\x{hex}{hex}),
// or invalid escape after reverse solidus (\{esc})
var ErrInvalidFormat = errors.New(("Invalid format"))
