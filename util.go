package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// TODO: Enable to control by option
var logFlag = false

func logPrintln(v ...interface{}) {
	if logFlag {
		log.Println(v)
	}
}

func markdownFileToHTML(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logPrintln(err)
		panic(err)
	}

	return blackfriday.Run(bytes, blackfriday.WithExtensions(blackfriday.CommonExtensions))
}
func runEditor(filename string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		fmt.Fprintf(os.Stderr, "Set $EDITOR\n")
		os.Exit(1)
	}

	splitted := strings.Split(editor, " ")
	logPrintln("splitted: %#v\n", splitted)
	cname := splitted[0]
	args := splitted[1:]
	args = append(args, filename)

	cmd := exec.Command(cname, args[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "editor error: %v\n", err)
		os.Exit(1)
	}
}
