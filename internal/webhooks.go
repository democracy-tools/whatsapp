package internal

import (
	"encoding/json"
	"net/http"

	"github.com/democracy-tools/whatsapp/internal/env"
	"github.com/sirupsen/logrus"
)

type Handle struct {
	webhookVerificationToken string
}

func NewHandle() *Handle {

	return &Handle{
		webhookVerificationToken: env.GetWhatsAppWebhookVerificationToken(),
	}
}

func (h *Handle) WebhookVerificationHandler(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query()
	mode := key.Get("hub.mode")
	token := key.Get("hub.verify_token")
	challenge := key.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == h.webhookVerificationToken {
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