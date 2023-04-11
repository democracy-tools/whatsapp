package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/democracy-tools/whatsapp/internal/env"
	"github.com/sirupsen/logrus"
)

type Api struct {
	webhookVerificationToken string
	slackUrl                 string
}

func NewApi() *Api {

	return &Api{
		webhookVerificationToken: env.GetWhatsAppWebhookVerificationToken(),
		slackUrl:                 env.GetSlackUrl(),
	}
}

func (api *Api) WebhookVerificationHandler(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query()
	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == api.webhookVerificationToken {
			w.Write([]byte(challenge))
			return
		}
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (api *Api) WebhookEventHandler(w http.ResponseWriter, r *http.Request) {

	var payload WebhookMessage
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logrus.Infof("failed to decode webhook message with '%v'")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pretty, err := buildMessage(payload)
	if err != nil {
		logrus.Infof("failed to marshal webhook message with '%v'")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = SendSlackMessage(api.slackUrl, string(pretty))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func buildMessage(message WebhookMessage) ([]byte, error) {

	if len(message.Entry) == 1 && len(message.Entry[0].Changes) == 1 {
		var res bytes.Buffer
		change := message.Entry[0].Changes[0]
		contact := change.Value.Contacts[0]
		res.WriteString(fmt.Sprintf("%s (%s)\n", contact.Profile.Name, contact.WaID))
		for _, currMessage := range change.Value.Messages {
			if currMessage.Type == "text" {
				res.WriteString(fmt.Sprintf("%s\n", currMessage.Text.Body))
			} else {
				res.WriteString(fmt.Sprintf("%s\n", currMessage.Type))
			}
		}
		return res.Bytes(), nil
	}

	return json.MarshalIndent(message, "", "  ")
}
