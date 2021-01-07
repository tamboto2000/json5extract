package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(
		`
{
	// this is single line comment	
	$id_num: 123,
	"msg": "hello world!",

	/* 
	this is
	multi line comment
	*/

	name: 'Franklin Collin Tamboto',
	'email': "tamboto2000@gmail.com",
	data: [123, 0x456, 'meta'],
}

["I'm an array!", "this is a string", 123, .456, 'I love you all <3']

"I'm a lonely string..."
'Don\'t worry, I\'m here for you'
"\"Me too!\"\
 Together ape string"

// some numbers
1234
0x1234
.90
123e45
Infinity
NaN
69.420 // nice 👌
666.666e+6666 // HAIL THE GARBAGE COLLECTOR
`)

	jsons, err := json5extract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	// iterate trough extracted result, determine the type, and print the parsed value
	for _, json := range jsons {
		// object
		if json.Kind == json5extract.Object {
			objMap, _ := json.Object()
			fmt.Println("Object values:")
			for key, val := range objMap {
				// string
				if val.Kind == json5extract.String {
					str, _ := val.String()
					fmt.Println("\t"+key+":", str)
				}

				// int
				if val.Kind == json5extract.Integer {
					i, _ := val.Integer()
					fmt.Println("\t"+key+":", i)
				}

				// array
				if val.Kind == json5extract.Array {
					fmt.Println("\t" + key + " (array):")
					vals, _ := val.Array()
					for _, val := range vals {
						// string
						if val.Kind == json5extract.String {
							str, _ := val.String()
							fmt.Println("\t\t", str)
						}

						// int
						if val.Kind == json5extract.Integer {
							i, _ := val.Integer()
							fmt.Println("\t\t", i)
						}

						// float
						if val.Kind == json5extract.Float {
							i, _ := val.Float()
							fmt.Println("\t\t", i)
						}
					}
				}
			}
		}

		// array
		if json.Kind == json5extract.Array {
			fmt.Println("Array values:")
			vals, _ := json.Array()
			for _, val := range vals {
				// string
				if val.Kind == json5extract.String {
					str, _ := val.String()
					fmt.Println("\t", str)
				}

				// int
				if val.Kind == json5extract.Integer {
					i, _ := val.Integer()
					fmt.Println("\t", i)
				}

				// float
				if val.Kind == json5extract.Float {
					i, _ := val.Float()
					fmt.Println("\t", i)
				}
			}
		}

		// string
		if json.Kind == json5extract.String {
			fmt.Println("String value:")
			str, _ := json.String()
			fmt.Println(str)
		}

		// int
		if json.Kind == json5extract.Integer {
			fmt.Println("Int value:")
			i, _ := json.Integer()
			fmt.Println(i)
		}

		// float
		if json.Kind == json5extract.Float {
			fmt.Println("Float value:")
			i, _ := json.Float()
			fmt.Println(i)
		}

		// infinite
		if json.Kind == json5extract.Infinity {
			fmt.Println("Infinity value:")
			i, _ := json.Infinity()
			fmt.Println(i)
		}

		// NaN
		if json.Kind == json5extract.NaN {
			fmt.Println("NaN value:")
			i, _ := json.NaN()
			fmt.Println(i)
		}
	}
}
