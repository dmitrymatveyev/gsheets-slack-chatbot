package model

// Message for posting to Slack
type Message struct {
	Text    string `json:"text"`
	Channel string `json:"channel"`
}
