package env

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetWhatsAppWebhookVerificationToken() string {

	return failIfEmpty("WHATSAPP_VERIFICATION_TOKEN")
}

func failIfEmpty(key string) string {

	res := os.Getenv(key)
	if res == "" {
		log.Fatalf("Please, add environment variable '%s'", key)
	}
	log.Debugf("%s: %s", key, res)

	return res
}
