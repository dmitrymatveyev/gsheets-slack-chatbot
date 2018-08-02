package main

import (
	"bytes"
	"gsheets-slack-chatbot/utility"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.
		Methods("POST").
		Path("/").
		HandlerFunc(post)

	log.Fatal(http.ListenAndServe(":80", r))
}

func post(w http.ResponseWriter, r *http.Request) {
	raw, _ := ioutil.ReadAll(r.Body)
	log.Printf("Received: %s\n", string(raw))

	var e OuterEvent
	err := utility.DeserializeJson(bytes.NewReader(raw), &e)
	if err != nil {
		utility.WriteBadRequest(w, err)
		return
	}

	if e.Type == "event_callback" {
		processMesChans(w, r, e.Event)
		utility.WriteResponse(w, e)
		return
	}

	if e.Type == "url_verification" {
		utility.WriteResponse(w, e)
		return
	}
}
