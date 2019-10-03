package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

var logFlag = false

func logPrintln(v ...interface{}) {
	if logFlag {
		log.Println(v)
	}
}

func runEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	logPrintln("$EDITOR", editor)

	if editor == "" {
		return errors.New("Set $EDITOR")
	}

	splitted := strings.Split(editor, " ")
	logPrintln("splitted: %#v\n", splitted)
	cname := splitted[0]
	args := splitted[1:]
	args = append(args, filename)

	cmd := exec.Command(cname, args[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
