# JSON5Extract

[![Go Reference](https://pkg.go.dev/badge/github.com/tamboto2000/json5extract.svg)](https://pkg.go.dev/github.com/tamboto2000/json5extract)

JSON5Extract is a library for extracting any valid JSON5 strings from a provided source like text, a string, bytes, or an an io.Reader written with Go from scratch. Extract JSON5 strings with ease without the hassle of regexp and other black magicks!

  - Extract JSON5 from []byte, string, file, or io.Reader
  - Access parsed value directly without unmarshal first
  - Backward compatible with JSON

For unmarshaling the bytes, you can use [this library](https://pkg.go.dev/github.com/flynn/json5)

### Installation

JSON5Extract requires Go v1.14 or up
Using Go Module is recommended
```sh
$ GO111MODULE=on go get github.com/tamboto2000/json5extract
```

# Examples

### Extract From String

```go
package main

import (
	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := `{
		$id_num: 123,
		"msg": "hello world!",
		name: 'Franklin Collin Tamboto',
		'email': "tamboto2000@gmail.com",
		data: [123, 0x456, 'meta'],
	}
	
	["I'm an array!", "this is a string", 123, .456, 'I love you all <3']`

	jsons, err := json5extract.FromString(raw)
	if err != nil {
		panic(err.Error())
	}

	// save result to ./extracted_jsons.json5
	if err := json5extract.Save(jsons); err != nil {
		panic(err.Error())
	}
}
```

### Extract From Bytes

```go
package main

import (
	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(`{
		$id_num: 123,
		"msg": "hello world!",
		name: 'Franklin Collin Tamboto',
		'email': "tamboto2000@gmail.com",
		data: [123, 0x456, 'meta'],
	}
	
	["I'm an array!", "this is a string", 123, .456, 'I love you all <3']`)

	jsons, err := json5extract.FromBytes(raw)
	if err != nil {
		panic(err.Error())
	}

	// save result
	if err := json5extract.Save(jsons); err != nil {
		panic(err.Error())
	}
}
```

### Extract From Reader
```go
package main

import (
	"bytes"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(`{
		$id_num: 123,
		"msg": "hello world!",
		name: 'Franklin Collin Tamboto',
		'email': "tamboto2000@gmail.com",
		data: [123, 0x456, 'meta'],
	}
	
	["I'm an array!", "this is a string", 123, .456, 'I love you all <3']`)

	r := bytes.NewReader(raw)

	jsons, err := json5extract.FromReader(r)
	if err != nil {
		panic(err.Error())
	}

	// save result
	if err := json5extract.Save(jsons); err != nil {
		panic(err.Error())
	}
}
```
### Determine Type and Retrieve value
For determining and accessing value, you need to check JSON5.Kind.
The Following code will demonstrate how to determine and retrieve parsed value

```go
package main

import (
	"fmt"

	"github.com/tamboto2000/json5extract"
)

func main() {
	raw := []byte(
		`
{
	$id_num: 123,
	"msg": "hello world!",
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
```
The output will be like this

```
Object values:
        data (array):
                 123
                 1110
                 meta
        $id_num: 123
        msg: hello world!
        name: Franklin Collin Tamboto
        email: tamboto2000@gmail.com 
Array values:
         I'm an array!
         this is a string
         123
         0.456
         I love you all <3
String value:
I'm a lonely string...
String value:
Don't worry, I'm here for you
String value:
"Me too!" Together ape string
Int value:
1234
Int value:
4660
Float value:
0.9
Float value:
1.23e+47
Infinity value:
+Inf
NaN value:
NaN
Float value:
69.42
Float value:
+Inf
```

### Retrieve Raw JSON5

To retrieve the extracted raw JSON5, use ```JSON5.RawBytes()``` for ```[]byte``` and ```JSON5.RawRunes()``` for ```[]rune```

### Identifier Standard For Object Identifier Name

This library follow strict mode for object identifier name, see [https://www.ecma-international.org/ecma-262/5.1/#sec-7.6](https://www.ecma-international.org/ecma-262/5.1/#sec-7.6) on "Identifier Names and Identifiers"

License
----

MIT
