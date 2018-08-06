package model

// MessageChannels holds infromation about message.channels event
type MessageChannels struct {
	Type        string `json:"type"`
	Channel     string `json:"channel"`
	User        string `json:"user"`
	Text        string `json:"text"`
	Ts          string `json:"ts"`
	EventTs     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`
}
