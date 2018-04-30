package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	flags "github.com/jessevdk/go-flags"
	"github.com/knadh/jsonconfig"
	"github.com/think-it-labs/actionizer/datastore"
	"github.com/think-it-labs/actionizer/models"
	"github.com/think-it-labs/actionizer/notifier"
	"github.com/think-it-labs/actionizer/server"
	"github.com/think-it-labs/actionizer/utils"
)

type Configuration struct {
	DataStorageConfig datastore.Config     `json:"datastore"`
	ActionDuration    utils.Duration       `json:"action_duration"`
	HTTPListen        string               `json:"http_listen"`
	HTTPPort          int                  `json:"http_port"`
	Slack             notifier.SlackConfig `json:"slack"`
}

type Options struct {
	Config  string `short:"c" long:"config" description:"Configuration file"`
	Verbose []bool `short:"v" long:"verbose" description:"Verbose"`
}

type ActionType int

const (
	NormalAction   ActionType = 0
	EnforcedAction ActionType = 1
)

type ActionRequest struct {
	Type  ActionType
	Users []models.User
	When  time.Time
}

func main() {

	// parse cli options
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(2)
	}

	if len(opts.Verbose) > 0 {
		log.SetLevel(log.DebugLevel)
	}

	configFile := opts.Config
	if configFile == "" {
		configFile = "actionizer.json"
	}

	// parse and load json config
	var config Configuration
	err = jsonconfig.Load(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v\n", err)
	}

	ds, err := datastore.New(config.DataStorageConfig)
	if err != nil {
		log.Fatalf("Datastorage error: %v\n", err)
	}

	slackChan := notifier.NewSlackNotifier(config.Slack)

	requestChan := make(chan ActionRequest)

	go func() {
		for actionRequest := range requestChan {
			var tasks []models.Task
			if actionRequest.Type == NormalAction {
				task, err := datastore.NewRandomTask(ds, actionRequest.When)
				if err != nil {
					log.Error(err)
				}
				tasks = append(tasks, task)
			} else if actionRequest.Type == EnforcedAction {
				for _, user := range actionRequest.Users {
					task, err := datastore.NewRandomEnforcedTask(ds, user, actionRequest.When)
					if err != nil {
						log.Error(err)
					}
					tasks = append(tasks, task)
				}
			}

			for _, task := range tasks {
				slackChan <- task
			}
		}
	}()

	// Weekly action picker
	go func() {

		// Get all non enforced tasks
		tasks := datastore.FilterTask(ds.GetCurrentTasks(), func(a models.Task) bool {
			return !a.Enforced
		})

		// Not weekly task found, create a new one.
		if len(tasks) == 0 {
			log.Infof("No weekly task found, creating new one.")
			// utils.StartOfWeek(time.Now())
			datastore.NewRandomTask(ds, utils.StartOfWeek(time.Now()))
		}

		for {
			nextWeek := utils.NextWeekStart(time.Now())
			sleepDuration := time.Until(nextWeek)
			log.Debugf("Waiting for weekend. Sleeping for %v, %v", sleepDuration, nextWeek)
			time.Sleep(sleepDuration)
			requestChan <- ActionRequest{
				Type: NormalAction,
				When: time.Now(),
			}
		}
	}()

	server := server.Server{
		Host: config.HTTPListen,
		Port: config.HTTPPort,
		DS:   ds,
		EnforceHandler: func(users []models.User, when time.Time) error {
			when = utils.StartOfWeek(when)
			requestChan <- ActionRequest{
				Type:  EnforcedAction,
				Users: users,
				When:  when,
			}
			return nil
		},
	}

	log.Printf("Listening on http://%s:%d", config.HTTPListen, config.HTTPPort)
	log.Fatal(server.Run())
}
