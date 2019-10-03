package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/browser"
	"github.com/rakyll/statik/fs"
	"github.com/urfave/cli"

	_ "github.com/sachaos/note/statik"
)

//go:generate statik -f -src=assets

func before(c *cli.Context) error {
	if c.GlobalBool("debug") {
		logFlag = true
	}
	return nil
}

func inject(path string, handler http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logPrintln("path: %s", r.URL.Path)

			if r.URL.Path == path {
				logPrintln("serve index")
				handler.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func render(fs http.FileSystem, filename string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := fs.Open(filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "file not found")
			return
		}
		_, err = io.Copy(w, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
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

	dirName := path.Dir(filename)
	logPrintln("dirname:", dirName)

	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	markupHandler := newMarkupServer(filename)
	go markupHandler.Start()
	defer markupHandler.Close()

	r := chi.NewRouter()
	r.Use(inject("/", render(statikFS, "/index.html")))

	r.Handle("/note-static/bundle.js", render(statikFS, "/note-static/bundle.js"))
	r.Handle("/ws", markupHandler)
	r.Handle("/*", http.FileServer(http.Dir(dirName)))

	go func() {
		logPrintln("Call http.ListenAndServe")
		if err := http.ListenAndServe(":1129", r); err != nil {
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
	app.Version = "0.4.0"

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
