package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// An email interface to send email (e.g. alert email)
type Email interface {
	Send(url string, secret string, email string, subject string, content string) error
}

// A struct that sends email using SendGrid service.
type SendGridEmail struct {
}

func NewSendGrid() Email {
	return &SendGridEmail{}
}

func (s *SendGridEmail) Send(url string, secret string, email string, subject string, content string) error {
	type Email struct {
		Email string `json:"email"`
	}

	type ContentElement struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type To struct {
		To []Email `json:"to"`
	}

	type SendGrid struct {
		Personalizations []To             `json:"personalizations"`
		From             Email            `json:"from"`
		Subject          string           `json:"subject"`
		Content          []ContentElement `json:"content"`
	}

	value := SendGrid{
		Personalizations: []To{
			{
				To: []Email{
					{
						Email: email,
					},
				},
			},
		},
		From: Email{
			Email: email,
		},
		Subject: subject,
		Content: []ContentElement{
			{
				Type:  "text/plain",
				Value: content,
			},
		},
	}

	json_data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	var client http.Client
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secret))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
