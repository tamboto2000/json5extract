package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	test := [][]byte{
		[]byte("1234"),
		[]byte("-1234"),
		[]byte("0xffff"),
		[]byte("-0xffff"),
		// this is invalid
		// []byte("-0xzz"),
		[]byte("0.0e4"),
		// only num 4 will be detected as valid
		// []byte("0.0ee4"),
		[]byte("0.0e-3"),
		[]byte("0.0e+3"),
		[]byte("0e+3"),
		[]byte("-.08"),
		[]byte("123e4"),
		[]byte("123e+4"),
		[]byte("123e-4"),
		// only num 8 will be detected as valid
		// []byte("08"),
		[]byte("123."),
		[]byte("Infinity"),
		[]byte("-Infinity"),
		[]byte("+Infinity"),
		[]byte("NaN"),
		[]byte("-NaN"),
		[]byte("+NaN"),
		// only num 45 will be detected as valid
		// []byte("e45"),
	}

	for _, byts := range test {
		jsons, err := json5extract.FromBytes(byts)
		if err != nil {
			fmt.Println(err.Error())
		}

		for _, json := range jsons {
			if json.Kind == json5extract.Integer {
				fmt.Println(json.Integer())
			}

			if json.Kind == json5extract.Float {
				fmt.Println(json.Float())
			}

			if json.Kind == json5extract.Infinity {
				fmt.Println(json.Infinity())
			}

			if json.Kind == json5extract.NaN {
				fmt.Println(json.NaN())
			}
		}
	}
}
