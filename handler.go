package main

import (
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type markupHandler struct {
	initialFileName string
	fw              *fsnotify.Watcher
	subscribers     []*websocket.Conn
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
	if err = ws.WriteMessage(websocket.TextMessage, result); err != nil {
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
