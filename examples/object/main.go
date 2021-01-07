package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(`{
		$id_num: 123,
		"msg": "hello world!",
		name: 'Franklin Collin Tamboto',
		'email': "tamboto2000@gmail.com",
		data: [123, 0x456, 'meta'],
	}`)

	jsons, err := json5extract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	for _, json := range jsons {
		if json.Kind == json5extract.Object {
			fmt.Println(string(json.RawBytes()))
		}
	}
}
