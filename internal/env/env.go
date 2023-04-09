package env

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetWhatAppToken() string {

	return failIfEmpty("WHATSAPP_TOKEN")
}

func GetWhatsAppFromPhone() string {

	return failIfEmpty("WHATSAPP_FROM_PHONE")
}

func GetWhatsAppWebhookVerificationToken() string {

	return failIfEmpty("WHATSAPP_VERIFICATION_TOKEN")
}

func GetSlackUrl() string {

	return failIfEmpty("SLACK_URL")
}

func failIfEmpty(key string) string {

	res := os.Getenv(key)
	if res == "" {
		log.Fatalf("Please, add environment variable '%s'", key)
	}
	log.Debugf("%s: %s", key, res)

	return res
}
