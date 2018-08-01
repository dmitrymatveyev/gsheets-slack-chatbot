package main

import (
	"encoding/json"
	"gsheets-slack-chatbot/utility"
	"io"
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
	e, err := getEventFromBody(r.Body)
	if err != nil {
		utility.WriteBadRequest(w, err)
		return
	}

	if e.Type == "url_verification" {
		utility.WriteResponse(w, e)
		return
	}
}

func getEventFromBody(body io.Reader) (Event, error) {
	raw, _ := ioutil.ReadAll(body)
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return Event{}, err
	}
	return e, nil
}
