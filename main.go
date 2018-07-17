package main

import (
	"fmt"
	"io/ioutil"
	"os"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

func main() {
	var bytes []byte
	var err error
	if len(os.Args) < 2 {
		bytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
	} else {
		filename := os.Args[1]
		bytes, err = ioutil.ReadFile(filename)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
	}
	result := blackfriday.Run(bytes, blackfriday.WithNoExtensions())
	fmt.Println(string(result))
}
