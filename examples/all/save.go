package main

import "os"

func saveBytes(data []byte, path string) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	if _, err := f.Write(data); err != nil {
		return err
	}

	return f.Close()
}

func saveString(data, path string) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return f.Close()
}
