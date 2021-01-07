package main

import (
	"github.com/tamboto2000/json5extract"
)

func main() {
	json5s, err := json5extract.FromFile("test.txt")
	if err != nil {
		panic(err.Error())
	}

	// there's only one JSON5 type extracted
	json5 := json5s[0]
	str, _ := json5.String()
	if err := saveString(str, "parsed.txt"); err != nil {
		panic(err.Error())
	}

	if err := saveBytes(json5.RawBytes(), "extracted.json5"); err != nil {
		panic(err.Error())
	}
}
