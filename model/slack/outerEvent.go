package model

import "encoding/json"

// OuterEvent holds basic information about an event
type OuterEvent struct {
	Token     string           `json:"token"`
	Challenge string           `json:"challenge"`
	Type      string           `json:"type"`
	APIAppID  string           `json:"api_app_id"`
	Event     *json.RawMessage `json:"event"`
}
