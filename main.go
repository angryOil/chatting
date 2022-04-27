package main

import (
	"encoding/json"
	"fmt"
	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
	"time"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}
func sendMessage(name, msg string) {
	msgCh <- Message{name, msg}
}

func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("add user: %s", username))
}
func deleteUsers(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("left user: %s", username))
}
func main() {
	msgCh = make(chan Message)
	es := eventsource.New(nil, nil)
	defer es.Close()

	go processMsgCh(es)

	var mux = pat.New()
	n := negroni.Classic()
	mux.Post("/message", postMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", deleteUsers)
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
