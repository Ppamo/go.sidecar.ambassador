package main

type Request struct {
	Context		RequestContext		`json:"context"`
	Data		map[string]interface{}  `json:"data"`
}

type RequestContext struct {
	Timestamp       string                  `json:"timestamp"`
	Application     string                  `json:"application"`
	ChannelId       string                  `json:"channel_id"`
}

type Mock struct {
	Body		map[string]interface{}	`json:"body"`
}
