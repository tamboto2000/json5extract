package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := [][]byte{
		[]byte("false"),
		[]byte("true"),
		// this string contain "true"
		[]byte("sdfsdfstruedfgsdf"),
		// this tring contain "false"
		[]byte("dzfvsdgffalsedgsdgsd"),
	}

	for _, r := range raw {
		jsons, err := json5extract.FromBytes(r)
		if err != nil {
			panic(err.Error())
		}

		for _, json := range jsons {
			if json.Kind == json5extract.Boolean {
				fmt.Println(json.Boolean())
			}
		}
	}
}
