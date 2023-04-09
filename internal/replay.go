package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/democracy-tools/whatsapp/internal/env"
	"github.com/sirupsen/logrus"
)

type TextMessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             MessageText `json:"text"`
}

type MessageText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

type Handle struct {
	auth string
	from string
}

func NewHandle() *Handle {

	return &Handle{
		auth: fmt.Sprintf("Bearer %s", env.GetWhatAppToken()),
		from: env.GetWhatsAppFromPhone(),
	}
}

func (h *Handle) Reply(w http.ResponseWriter, r *http.Request) {

	to, message := r.URL.Query().Get("to"), r.URL.Query().Get("m")
	if to == "" {
		w.Write([]byte("Please, specify 'to' as query string param"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if message == "" {
		w.Write([]byte("Please, specify 'm' (message) as query string param"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	send(h.auth, h.from, to, message)
}

func send(auth string, from string, to string, message string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(TextMessageRequest{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "text",
		Text: MessageText{
			PreviewURL: false,
			Body:       message,
		},
	})
	if err != nil {
		logrus.Errorf("failed to encode whatsapp message request with '%v'. target phone '%s'", err, to)
		return err
	}

	r, err := http.NewRequest(http.MethodPost, getMessageUrl(from), &buf)
	if err != nil {
		logrus.Errorf("failed to create HTTP request for sending a whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", auth)

	client := http.Client{}
	response, err := client.Do(r)
	if err != nil {
		logrus.Errorf("failed to send whatsapp message to '%s' with '%v'", to, err)
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		logrus.Infof("failed to send whatsapp message to '%s' with '%s'", to, response.Status)
		return err
	}

	return nil
}

func getMessageUrl(from string) string {

	return fmt.Sprintf("https://graph.facebook.com/v16.0/%s/messages", from)
}
