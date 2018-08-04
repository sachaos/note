package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/rakyll/statik/fs"
	blackfriday "gopkg.in/russross/blackfriday.v2"

	_ "github.com/sachaos/md2html/statik"
)

// TODO: Enable to control by option
var logFlag = false

func markdownFileToHTML(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logPrintln(err)
		panic(err)
	}

	return blackfriday.Run(bytes, blackfriday.WithExtensions(blackfriday.CommonExtensions))
}

//go:generate statik -f -src=html

func logPrintln(v ...interface{}) {
	if logFlag {
		log.Println(v)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type mdToHTML struct {
	initialFileName string
	fw              *fsnotify.Watcher
	subscribers     []*websocket.Conn
}

// TODO: It should be return err
func newMdToHTML(filename string) *mdToHTML {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logPrintln(err)
		panic(err)
	}
	return &mdToHTML{initialFileName: filename, fw: watcher}
}

func (m *mdToHTML) Close() error {
	return m.fw.Close()
}

// TODO: It should be return err
func (m *mdToHTML) AddFile(filename string) {
	if err := m.fw.Add(filename); err != nil {
		logPrintln(err)
		panic(err)
	}
}

func (m *mdToHTML) Start() {
	for {
		select {
		case event := <-m.fw.Events:
			for _, ws := range m.subscribers {
				logPrintln("event:", event)
				result := markdownFileToHTML(event.Name)
				if err := ws.WriteMessage(websocket.TextMessage, result); err != nil {
					logPrintln(err)
					return
				}
			}
		case err := <-m.fw.Errors:
			logPrintln("WatchError:", err)
			return
		}
	}
}

func (m *mdToHTML) Subscribe(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logPrintln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	m.AddFile(m.initialFileName)

	result := markdownFileToHTML(m.initialFileName)
	if err = ws.WriteMessage(websocket.TextMessage, result); err != nil {
		logPrintln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.subscribers = append(m.subscribers, ws)

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
	ws.Close()
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

	mdToHTML := newMdToHTML(filename)
	go mdToHTML.Start()
	defer mdToHTML.Close()

	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	http.HandleFunc("/ws", mdToHTML.Subscribe)

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
