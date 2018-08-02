package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/log"
	"gsheets-slack-chatbot/model"
	"net/http"
)

const token = "xoxa-409706946007-408390633905-409665283381-1e30506391a0a7f939df0683e2fe7237"
const url = "https://slack.com/api/chat.postMessage"
const botUserID = "UC0BGJMSM"

type message struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}

// ProcessMessageChannels handles an event processing
func ProcessMessageChannels(mesChan *model.MessageChannels) {
	where := "processor.ProcessMessageChannels(...)"

	if mesChan.User == botUserID {
		log.Trace(where, "It's me. Skip.")
		return
	}

	mes := message{
		Channel: mesChan.Channel,
		Text:    mesChan.Text,
	}

	sendMessage(&mes)
}

func sendMessage(m *message) {
	where := "processor.sendMessage(...)"

	raw, err := json.Marshal(&m)
	if err != nil {
		log.Error(where, "Couldn't serialize message.", err)
	}
	log.Trace(where, fmt.Sprintf("Sucessfully serialized message: %s", string(raw)))

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		log.Error(where, "Couldn't create request object.", err)
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	var c http.Client
	log.Trace(where, fmt.Sprintf("Making call to %s.", url))
	res, err := c.Do(req)
	if err != nil {
		log.Error(where, fmt.Sprintf("Couldn't call to %s.", url), err)
	}
	log.Trace(where, fmt.Sprintf("Successfully called to %s. Status: %s", url, res.Status))
}
