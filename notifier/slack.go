package notifier

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
	"github.com/think-it-labs/actionizer/models"
)

var params = slack.PostMessageParameters{
	Username: "Actionizer",
	AsUser:   false,
	Markdown: true,
}

type SlackConfig struct {
	TokenEnv string `json:"token_env"`
	Channel  string `json:"channel"`
}

func NewSlackNotifier(config SlackConfig) chan<- models.Task {
	ch := make(chan models.Task, 32)
	api := slack.New(os.Getenv(config.TokenEnv))
	go func() {
		for task := range ch {
			continue
			text := fmt.Sprintf("*New <https://actionizer.think-it.io/|Actionizer> Item:*\n*%s:* _%s_", task.User.Fullname, task.Action.Description)
			_, _, err := api.PostMessage(config.Channel, text, params)
			if err != nil {
				log.Errorf("Error sending slack notification: %v", err)
			}

		}
	}()
	return ch
}
