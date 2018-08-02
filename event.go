package main

import (
	"encoding/json"
	"gsheets-slack-chatbot/utility"
	"log"
	"net/http"
)

const botUserID = "UC0BGJMSM"

type MessageChannels struct {
	Type        string `json:"type"`
	Channel     string `json:"channel"`
	User        string `json:"user"`
	Text        string `json:"text"`
	Ts          string `json:"ts"`
	EventTs     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`
}

func processMesChans(w http.ResponseWriter, r *http.Request, raw *json.RawMessage) {
	var e MessageChannels
	err := json.Unmarshal(*raw, &e)
	if err != nil {
		utility.WriteBadRequest(w, err)
		return
	}
	if e.User == botUserID {
		return
	}
	log.Printf("Received: %s\n", e.Text)

	go process(&e)
}
