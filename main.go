package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"github.com/rakyll/statik/fs"
	_ "github.com/sachaos/mu/statik"
)

//go:generate statik -f -src=html

func main() {
	var err error

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Specify filename by argument")
		os.Exit(1)
	}

	filename := os.Args[1]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			logPrintln(err)
			os.Exit(1)
		}
		if err = file.Close(); err != nil {
			logPrintln(err)
			os.Exit(1)
		}
	}

	statikFS, err := fs.New()
	if err != nil {
		logPrintln(err)
		panic(err)
	}

	markupHandler := newMarkupServer(filename)
	go markupHandler.Start()
	defer markupHandler.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	http.Handle("/ws", markupHandler)

	go func() {
		if err = http.ListenAndServe(":1129", nil); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err = browser.OpenURL("http://localhost:1129"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}()

	runEditor(filename)
}
