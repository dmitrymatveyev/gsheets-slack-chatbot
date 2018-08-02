package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const token = "xoxa-409706946007-408390633905-409665283381-1e30506391a0a7f939df0683e2fe7237"
const url = "https://slack.com/api/chat.postMessage"

type Message struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

func process(e *MessageChannels) {
	m := Message{
		Channel: e.Channel,
		Text:    e.Text,
	}

	raw, err := json.Marshal(&m)
	if err != nil {
		log.Println(err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	var c http.Client
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println(res.StatusCode)
}
