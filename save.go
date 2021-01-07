package json5extract

import "os"

// SaveToPath save extracted JSONs to a file path
func SaveToPath(data []*JSON5, path string) error {
	return save(data, path)
}

// Save save extracted JSONs to ./extracted_jsons.json5
func Save(data []*JSON5) error {
	return save(data, "./extracted_jsons.json5")
}

func save(data []*JSON5, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	rest := make([]byte, 0)
	rest = append(rest, 91)
	c := len(data)
	if c > 0 {
		for i, d := range data {
			rest = append(rest, d.RawBytes()...)
			if i == c-1 {
				rest = append(rest, 93)
			} else {
				rest = append(rest, 44)
			}
		}
	} else {
		rest = append(rest, 93)
	}

	if _, err = f.Write(rest); err != nil {
		return err
	}

	return nil
}
