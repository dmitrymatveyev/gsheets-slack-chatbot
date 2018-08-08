package main

import (
	proc "gsheets-slack-chatbot/processor"
	util "gsheets-slack-chatbot/utility"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	where := "main.main()"

	config, err := util.NewConfig("config.json")
	if err != nil {
		panic(err)
	}
	log, err := util.NewLog(config)
	if err != nil {
		panic(err)
	}
	defer log.Close()
	helper, err := util.NewWebHelper()
	if err != nil {
		panic(err)
	}
	proc, err := proc.New(log, config)
	if err != nil {
		panic(err)
	}
	contr, err := NewController(log, helper, proc)
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.
		Methods("POST").
		Path("/").
		HandlerFunc(contr.Post)

	log.Trace(where, "Listening on port 80.")
	err = http.ListenAndServe(":80", router)

	log.Error(where, "", err)
	os.Exit(1)
}
