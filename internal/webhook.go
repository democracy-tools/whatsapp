package internal

import (
	"encoding/json"
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

	pretty, err := json.MarshalIndent(payload, "", "  ")
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
