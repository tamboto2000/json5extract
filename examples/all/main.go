package main

import (
	"github.com/tamboto2000/json5extract"
)

func main() {
	json5s, err := json5extract.FromFile("test_heavy.txt")
	if err != nil {
		panic(err.Error())
	}

	// filtered := make([]*json5extract.JSON5, 0)
	// for _, json := range json5s {
	// 	if json.Kind == json5extract.String {
	// 		filtered = append(filtered, json)
	// 	}
	// }

	if err := json5extract.Save(json5s); err != nil {
		panic(err.Error())
	}
}
