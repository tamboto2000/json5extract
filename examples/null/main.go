package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := [][]byte{
		[]byte("null"),
		// this string contains "null"
		[]byte("sdfsasfnullsdfasfas"),
	}

	for _, r := range raw {
		jsons, err := json5extract.FromBytes(r)
		if err != nil {
			panic(err.Error())
		}

		for _, json := range jsons {
			if json.Kind == json5extract.Null {
				// null := json.(*json5extract.Null)
				fmt.Println("well, it just null...")
			}
		}
	}
}
