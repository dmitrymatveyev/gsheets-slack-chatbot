package main

import (
	"encoding/json"
	"fmt"
	"gsheets-slack-chatbot/model"
	proc "gsheets-slack-chatbot/processor"
	"gsheets-slack-chatbot/utility"
	"io/ioutil"
	"net/http"
	"reflect"
)

// Controller handles requests
type Controller struct {
	log    *utility.Log
	helper *utility.WebHelper
	proc   *proc.Processor
}

// NewController creates new instance of Controller
func NewController(log *utility.Log, helper *utility.WebHelper, proc *proc.Processor) (*Controller, error) {
	return &Controller{log: log, helper: helper, proc: proc}, nil
}

// Post handles POST request to root
func (c *Controller) Post(w http.ResponseWriter, r *http.Request) {
	where := "main.Controller.Post(...)"

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.log.Error(where, "Failed to read a body.", err)
		c.helper.WriteBadRequest(w, err)
		return
	}
	c.log.Trace(where, fmt.Sprintf("Received: \"%s\".", string(raw)))

	var e model.OuterEvent
	err = json.Unmarshal(raw, &e)
	if err != nil {
		c.log.Error(where, fmt.Sprintf("Failed to deserialize %v.", reflect.TypeOf(e)), err)
		c.helper.WriteBadRequest(w, err)
		return
	}
	c.log.Trace(where, fmt.Sprintf("Successfully deserialized %v.", reflect.TypeOf(e)))

	if e.Type == "event_callback" {
		c.handleMessage(w, r, e.Event)
		return
	}

	if e.Type == "url_verification" {
		c.helper.WriteResponse(w, e)
		return
	}
}

func (c *Controller) handleMessage(w http.ResponseWriter, r *http.Request, raw *json.RawMessage) {
	where := "main.Controller.handleMessage(...)"

	var m model.MessageChannels
	err := json.Unmarshal(*raw, &m)
	if err != nil {
		c.log.Error(where, fmt.Sprintf("Failed to deserialize %v.", reflect.TypeOf(m)), err)
		c.helper.WriteBadRequest(w, err)
		return
	}
	c.log.Trace(where, fmt.Sprintf("Successfully deserialized %v.", reflect.TypeOf(m)))

	go c.proc.ProcessMessageChannels(&m)

	w.WriteHeader(http.StatusOK)
}
