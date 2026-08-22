package domain

import "context"

const EventSend = "sms.send"

type SendRequest struct {
	To   string `json:"to"`
	Text string `json:"text"`
}

type Provider interface {
	Name() string
	Send(ctx context.Context, to, text string) error
}
