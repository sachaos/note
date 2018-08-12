package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type message struct {
	Title string `json:"title"`
	HTML  string `json:"html"`
	Lines []int  `json:"lines"`
}

func createMessage(html []byte, title string, lines []int) []byte {
	msg := message{
		Title: title,
		HTML:  injectAnchor(html),
		Lines: lines,
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		logPrintln(err)
		return []byte{}
	}
	return bytes
}

func injectAnchor(html []byte) string {
	var out string
	for i, str := range strings.Split(string(html), "\n") {
		out += fmt.Sprintf("<a id=\"%d\"></a>", i) + str + "\n"
	}
	return out
}

type markupHandler struct {
	initialFileName string
	fw              *fsnotify.Watcher
	subscribers     []*websocket.Conn
	current         []byte
}

// TODO: It should be return err
func newMarkupServer(filename string) *markupHandler {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logPrintln(err)
		panic(err)
	}
	return &markupHandler{initialFileName: filename, fw: watcher}
}

func (m *markupHandler) Close() error {
	return m.fw.Close()
}

// TODO: It should be return err
func (m *markupHandler) AddFile(filename string) {
	if err := m.fw.Add(filename); err != nil {
		logPrintln(err)
		panic(err)
	}
}

func (m *markupHandler) Start() {
	for {
		select {
		case event := <-m.fw.Events:
			for _, ws := range m.subscribers {
				logPrintln("event:", event)
				result := markdownFileToHTML(event.Name)
				lines := diffLines(string(m.current), string(result))
				m.current = result
				msg := createMessage(result, m.initialFileName, lines)

				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					logPrintln("WriteMessage Error:", err)
					return
				}
			}
		case err := <-m.fw.Errors:
			logPrintln("Watch Error:", err)
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (m *markupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	m.current = result
	msg := createMessage(result, m.initialFileName, []int{})
	if err = ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		logPrintln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m.subscribers = append(m.subscribers, ws)

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			logPrintln(err)
			break
		}
	}
	ws.Close()
}

func markdownFileToHTML(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logPrintln(err)
		panic(err)
	}

	return blackfriday.Run(bytes, blackfriday.WithExtensions(blackfriday.CommonExtensions))
}
