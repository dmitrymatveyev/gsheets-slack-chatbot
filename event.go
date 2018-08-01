package main

type Event struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}
