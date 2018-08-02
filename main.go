package main

import (
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/log"
	"gsheets-slack-chatbot/model"
	"gsheets-slack-chatbot/processor"
	"gsheets-slack-chatbot/web"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("POST").
		Path("/").
		HandlerFunc(post)

	err := http.ListenAndServe(":80", router)

	log.Error("main.main()", "", err)
	os.Exit(1)
}

func post(w http.ResponseWriter, r *http.Request) {
	where := "main.post(...)"

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(where, "Failed to read a body.", err)
		web.WriteBadRequest(w, err)
		return
	}
	log.Trace(where, fmt.Sprintf("Received: \"%s\".", string(raw)))

	var e model.OuterEvent
	err = json.Unmarshal(raw, &e)
	if err != nil {
		log.Error(where, fmt.Sprintf("Failed to deserialize %v.", reflect.TypeOf(e)), err)
		web.WriteBadRequest(w, err)
		return
	}
	log.Trace(where, fmt.Sprintf("Successfully deserialized %v.", reflect.TypeOf(e)))

	if e.Type == "event_callback" {
		handleMessage(w, r, e.Event)
		return
	}

	if e.Type == "url_verification" {
		web.WriteResponse(w, e)
		return
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request, raw *json.RawMessage) {
	where := "main.handleMessage(...)"

	var m model.MessageChannels
	err := json.Unmarshal(*raw, &m)
	if err != nil {
		log.Error(where, fmt.Sprintf("Failed to deserialize %v.", reflect.TypeOf(m)), err)
		web.WriteBadRequest(w, err)
		return
	}
	log.Trace(where, fmt.Sprintf("Successfully deserialized %v.", reflect.TypeOf(m)))

	go processor.ProcessMessageChannels(&m)

	w.WriteHeader(http.StatusOK)
}
