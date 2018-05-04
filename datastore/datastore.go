package datastore

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/think-it-labs/actionizer/models"
	"github.com/think-it-labs/actionizer/utils"
)

type DataStore interface {
	AddUser(models.User) error
	AddTask(models.Task) error
	AddAction(models.Action) error

	GetUsers() (map[string]models.User, error)
	GetActions() (map[string]models.Action, error)
	GetTasks() (map[string]models.Task, error)

	GetCurrentTasks() map[string]models.Task
}

type Config struct {
	Backend string            `json:"backend"`
	Params  map[string]string `json:"params"`
}

type CreatorFunc func(Config) (DataStore, error)

var backends map[string]CreatorFunc

// temporarly solution. TODO: fix
type ActionModifierFunc func(DataStore, *models.Action) *models.Action

var actionsModifiers map[string]ActionModifierFunc

var requireRandUserList = []string{
	"Be the right-hand of another Think.iteer for 1 hour",
	"Jynx: Don't Talk untill someone  calls your name",
}

var requireRandColorList = []string{
	"Wear the same color for two days",
}

func init() {
	actionsModifiers = make(map[string]ActionModifierFunc)

	for _, specialAction := range requireRandUserList {
		actionsModifiers[specialAction] = func(ds DataStore, action *models.Action) *models.Action {
			users, err := ds.GetUsers()
			if err != nil {
				log.Errorf("Cannot load users")
				return action
			}

			choosenUser := pickRandomUser(users)
			action.Description = fmt.Sprintf(action.Description, choosenUser.Fullname)

			return action
		}
	}

	for _, specialAction := range requireRandColorList {
		actionsModifiers[specialAction] = func(ds DataStore, action *models.Action) *models.Action {
			action.Description = fmt.Sprintf(action.Description, utils.RandColor())
			return action
		}
	}

}

var backendNotSupportedError = errors.New("Storage backend not supported")

// RegisterBackend register a new data store backend
// Usually this method is called inside the `init` function of a backend
func RegisterBackend(backend string, creator CreatorFunc) {
	if backends == nil {
		backends = make(map[string]CreatorFunc)
	}
	backends[backend] = creator
}

func New(config Config) (DataStore, error) {
	if backends == nil {
		backends = make(map[string]CreatorFunc)
	}
	creator, ok := backends[config.Backend]
	if !ok {
		return nil, backendNotSupportedError
	}
	return creator(config)
}

func FilterAction(actions map[string]models.Action, checkFunc func(models.Action) bool) map[string]models.Action {
	filteredActions := make(map[string]models.Action)

	for key, action := range actions {
		if checkFunc(action) {
			filteredActions[key] = action
		}
	}
	return filteredActions
}

func FilterTask(tasks map[string]models.Task, checkFunc func(models.Task) bool) map[string]models.Task {
	filteredTasks := make(map[string]models.Task)

	for key, task := range tasks {
		if checkFunc(task) {
			filteredTasks[key] = task
		}
	}
	return filteredTasks
}

func NewRandomTask(ds DataStore, start time.Time) (models.Task, error) {
	rand.Seed(time.Now().UnixNano())

	users, err := ds.GetUsers()
	if err != nil {
		return models.Task{}, err
	}

	actions, err := ds.GetActions()
	if err != nil {
		return models.Task{}, err
	}

	choosenUser := pickRandomUser(users)
	log.Debugf("Choosen User: %+v", choosenUser)
	if choosenUser.Remote {
		actions = FilterAction(actions, func(a models.Action) bool {
			return a.Remotee
		})
	}

	choosenAction := pickRandomAction(actions)
	log.Debugf("Choosen Action: %+v", choosenAction)

	if modifierCall, ok := actionsModifiers[choosenAction.ID]; ok {
		choosenAction = modifierCall(ds, choosenAction)
		log.Debugf("After modification: %+v", choosenAction)
	}

	taskID := uuid.Must(uuid.NewV4())

	task := models.Task{
		ID:        taskID.String(),
		User:      *choosenUser,
		Action:    *choosenAction,
		Deadline:  start.Add(time.Duration(choosenAction.Duration)),
		StartDate: start,
	}

	err = ds.AddTask(task)
	return task, err
}

func NewRandomEnforcedTask(ds DataStore, user models.User, start time.Time) (models.Task, error) {
	rand.Seed(time.Now().UnixNano())

	actions, _ := ds.GetActions()
	actions = FilterAction(actions, func(a models.Action) bool {
		return a.Enforce
	})

	choosenAction := pickRandomAction(actions)
	log.Debugf("Choosen Action for user %q: %+v", user.Fullname, choosenAction)

	if modifierCall, ok := actionsModifiers[choosenAction.ID]; ok {
		choosenAction = modifierCall(ds, choosenAction)
		log.Debugf("After modification: %+v", choosenAction)
	}

	taskID := uuid.Must(uuid.NewV4())

	task := models.Task{
		ID:        taskID.String(),
		User:      user,
		Action:    *choosenAction,
		Deadline:  start.Add(time.Duration(choosenAction.Duration)),
		StartDate: start,
		Enforced:  true,
	}

	err := ds.AddTask(task)
	return task, err
}

func pickRandomUser(users map[string]models.User) *models.User {

	var usersNames []string
	for user := range users {
		usersNames = append(usersNames, user)
	}

	userName := usersNames[rand.Intn(len(usersNames))]
	user := users[userName]
	return &user
}

func pickRandomAction(actions map[string]models.Action) *models.Action {
	var actionsNames []string
	for action := range actions {
		actionsNames = append(actionsNames, action)
	}

	log.Infof("Len actions: %d\n", len(actionsNames))
	actionName := actionsNames[rand.Intn(len(actionsNames))]
	action := actions[actionName]
	return &action
}
