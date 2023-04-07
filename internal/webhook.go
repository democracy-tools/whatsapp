package internal

import (
	"encoding/json"
	"net/http"

	"github.com/democracy-tools/whatsapp/internal/env"
	"github.com/sirupsen/logrus"
)

type Api struct {
	webhookVerificationToken string
}

func NewApi() *Api {

	return &Api{
		webhookVerificationToken: env.GetWhatsAppWebhookVerificationToken(),
	}
}

func (api *Api) WebhookVerificationHandler(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query()
	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == api.webhookVerificationToken {
			w.WriteHeader(http.StatusOK)
			data, err := json.Marshal(challenge)
			if err != nil {
				logrus.Errorf("failed to marshal challenge with '%v'", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
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
	w.WriteHeader(http.StatusAccepted)
}
