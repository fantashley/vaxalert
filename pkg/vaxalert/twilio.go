package vaxalert

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kevinburke/twilio-go"
)

type Alerter interface {
	Alert(context.Context, string) error
}

type TwilioAlerter struct {
	messagingSid string
	toNumbers    []string
	client       *twilio.Client
}

func NewTwilioAlerter(
	accountSid,
	authToken,
	messagingSid string,
	httpClient *http.Client,
	toNumbers []string,
) TwilioAlerter {
	client := twilio.NewClient(accountSid, authToken, httpClient)
	return TwilioAlerter{
		messagingSid: messagingSid,
		toNumbers:    toNumbers,
		client:       client,
	}
}

func (a TwilioAlerter) Alert(ctx context.Context, message string) error {
	for _, num := range a.toNumbers {
		data := url.Values{
			"To":                  {num},
			"MessagingServiceSid": {a.messagingSid},
			"Body":                {message},
		}
		msg, err := a.client.Messages.Create(ctx, data)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}
		if msg == nil {
			return errors.New("created message is nil")
		}
	}

	return nil
}
