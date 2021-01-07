package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(`[
			123, 
			"123", 
			0x123, 
			["Hello World!", [123, "456", 78e9],   ],
			true,
			false,
			null,
			Infinity,
			NaN,
			1e23,
			1e-23,
		]`)

	jsons, err := json5extract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		fmt.Println(string(json.RawBytes()))
		fmt.Println(json.Array())
	}
}
