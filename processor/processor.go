package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/model"
	util "gsheets-slack-chatbot/utility"
	"net/http"
)

// Processor processes a message
type Processor struct {
	log    *util.Log
	config *util.Config
}

// New creates new instance of Processor
func New(log *util.Log, config *util.Config) (*Processor, error) {
	return &Processor{log: log, config: config}, nil
}

// ProcessMessageChannels handles an event processing
func (p *Processor) ProcessMessageChannels(mesChan *model.MessageChannels) {
	where := "processor.Processor.ProcessMessageChannels(...)"

	botUserID, err := p.config.Get("SlackBotUserID")
	if err != nil {
		p.log.Error(where, "", err)
		return
	}

	if mesChan.User == botUserID {
		p.log.Trace(where, "It's me. Skip.")
		return
	}

	mes := message{
		Channel: mesChan.Channel,
		Text:    mesChan.Text,
	}

	p.sendMessage(&mes)
}

func (p *Processor) sendMessage(m *message) {
	where := "processor.Processor.sendMessage(...)"

	raw, err := json.Marshal(&m)
	if err != nil {
		p.log.Error(where, "Couldn't serialize message.", err)
		return
	}
	p.log.Trace(where, fmt.Sprintf("Sucessfully serialized message: %s", string(raw)))

	url, err := p.config.Get("SlackPostURL")
	if err != nil {
		p.log.Error(where, "", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		p.log.Error(where, "Couldn't create request object.", err)
		return
	}

	token, err := p.config.Get("SlackWorkspaceToken")
	if err != nil {
		p.log.Error(where, "", err)
		return
	}

	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	var c http.Client
	p.log.Trace(where, fmt.Sprintf("Making call to %s.", url))
	res, err := c.Do(req)
	if err != nil {
		p.log.Error(where, fmt.Sprintf("Couldn't call to %s.", url), err)
		return
	}
	p.log.Trace(where, fmt.Sprintf("Successfully called to %s. Status: %s", url, res.Status))
}
