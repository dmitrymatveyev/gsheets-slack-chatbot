package main

import "encoding/json"

type OuterEvent struct {
	Token     string           `json:"token"`
	Challenge string           `json:"challenge"`
	Type      string           `json:"type"`
	APIAppID  string           `json:"api_app_id"`
	Event     *json.RawMessage `json:"event"`
}
