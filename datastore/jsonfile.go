package datastore

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/think-it-labs/actionizer/models"
	"github.com/think-it-labs/actionizer/utils"
)

type JSONStore struct {
	actions map[string]models.Action
	users   map[string]models.User
	tasks   map[string]models.Task

	actionsFile   string
	usersFile     string
	tasksFile     string
	doneTasksFile string
}

func init() {
	RegisterBackend("jsonfile", load)
}

func syncFiles(fileName string, target interface{}) error {

	syncFunc := func() error {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(target)
		if err != nil {
			log.Errorf("Cannot sync with file %s: %v", fileName, err)
		}
		return err
	}

	utils.AddFileWatcher(fileName, func() {
		syncFunc()
	})

	return syncFunc()

}

func load(config Config) (DataStore, error) {
	jsonStore := JSONStore{
		usersFile:     config.Params["users_file"],
		actionsFile:   config.Params["actions_file"],
		tasksFile:     config.Params["tasks_file"],
		doneTasksFile: config.Params["done_tasks_file"],
	}

	err := syncFiles(jsonStore.usersFile, &jsonStore.users)
	if err != nil {
		return nil, err
	}

	err = syncFiles(jsonStore.actionsFile, &jsonStore.actions)
	if err != nil {
		return nil, err
	}

	err = syncFiles(jsonStore.tasksFile, &jsonStore.tasks)
	if err != nil {
		jsonStore.tasks = make(map[string]models.Task)
	}

	log.Infof("%d user loaded", len(jsonStore.users))
	log.Infof("%d action loaded", len(jsonStore.actions))

	return &jsonStore, nil
}

func (js *JSONStore) GetUsers() (map[string]models.User, error) {
	return js.users, nil
}

func (js *JSONStore) GetActions() (map[string]models.Action, error) {
	return js.actions, nil
}

func (js *JSONStore) GetTasks() (map[string]models.Task, error) {
	return js.tasks, nil
}

func (js *JSONStore) AddUser(user models.User) error {
	js.users[user.Fullname] = user
	return persist(js.users, js.usersFile)
}

func (js *JSONStore) AddAction(action models.Action) error {
	js.actions[action.ID] = action
	return persist(js.actions, js.actionsFile)
}

func (js *JSONStore) AddTask(task models.Task) error {
	js.tasks[task.ID] = task
	return persist(js.tasks, js.tasksFile)
}

func (js *JSONStore) GetCurrentTasks() map[string]models.Task {
	tasks := make(map[string]models.Task)
	// now := time.Now()
	for _, task := range js.tasks {
		if task.Done {
			continue
		}

		tasks[task.ID] = task
	}
	return tasks
}

func persist(what interface{}, where string) error {
	file, err := os.Create(where)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(what)
}
