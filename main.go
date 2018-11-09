package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/rakyll/statik/fs"
	_ "github.com/sachaos/note/statik"
	"github.com/urfave/cli"
)

//go:generate statik -f -src=assets

func before(c *cli.Context) error {
	if c.GlobalBool("debug") {
		logFlag = true
	}
	return nil
}

func run(c *cli.Context) error {
	var err error

	logPrintln("c.Args(): ", c.Args())

	if len(c.Args()) != 1 {
		return errors.New("Specify filename by argument")
	}

	filename := c.Args()[0]
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}
	}

	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	markupHandler := newMarkupServer(filename)
	go markupHandler.Start()
	defer markupHandler.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	http.Handle("/ws", markupHandler)

	go func() {
		logPrintln("Call http.ListenAndServe")
		if err := http.ListenAndServe(":1129", nil); err != nil {
			logPrintln(err)
			os.Exit(1)
		}
	}()

	logPrintln("Call openURL")
	if err = browser.OpenURL("http://localhost:1129"); err != nil {
		return err
	}

	logPrintln("Call runEditor()")

	if c.Bool("no-editor") {
		for {
			time.Sleep(1 * time.Second)
		}
	} else {
		return runEditor(filename)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "note"
	app.Usage = "Realtime markdown previewer"
	app.Version = "0.2.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Hidden: true,
		},
		cli.BoolFlag{
			Name: "no-editor",
		},
	}

	app.Before = cli.BeforeFunc(before)
	app.Action = cli.ActionFunc(run)

	if err := app.Run(os.Args); err != nil {
		logPrintln(err)
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
