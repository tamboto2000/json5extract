package json5extract

import "os"

// FromFile extract JSON5 strings from a file in path
func FromFile(path string) ([]JSON5, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader, err := readFromReader(f)
	if err != nil {
		return nil, err
	}

	return parseAll(reader)
}

// FromBytes extract JSON5 strings from array of bytes
func FromBytes(byts []byte) ([]JSON5, error) {
	r := readFromBytes(byts)
	return parseAll(r)
}
