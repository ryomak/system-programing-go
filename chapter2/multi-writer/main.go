package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	var str string
	fmt.Scan(&str)
	file, err := os.Create("test.text")
	if err != nil {
		panic(err)
	}
	writers := io.MultiWriter(os.Stdout, file)
	writers.Write([]byte(str))
}
