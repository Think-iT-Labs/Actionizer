package datastore

import (
	"time"

	"github.com/think-it-labs/actionizer/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Database struct {
	mongodb *mgo.Database
}

type taskAssociation struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id,omitempty"`
	ActionID bson.ObjectId `json:"action_id" bson:"action_id,omitempty"`
	Deadline time.Time     `json:"deadline" bson:"deadline,omitempty"`
}

func init() {
	RegisterBackend("mongodb", connect)
}

func connect(config Config) (DataStore, error) {
	session, err := mgo.Dial(config.Params["Host"])
	if err != nil {
		return nil, err
	}

	db := session.DB(config.Params["Name"])

	// If User is not empty we do an auth
	if config.Params["User"] != "" {
		err := db.Login(config.Params["User"], config.Params["Password"])
		if err != nil {
			return nil, err
		}
	}
	return &Database{db}, nil
}

func (db *Database) ucol() *mgo.Collection {
	return db.mongodb.C("users")
}

func (db *Database) acol() *mgo.Collection {
	return db.mongodb.C("actions")
}

func (db *Database) tcol() *mgo.Collection {
	return db.mongodb.C("tasks")
}

func (db *Database) GetCurrentTask() (*models.Task, error) {
	var taskAssociationItem taskAssociation
	now := time.Now().UTC()
	err := db.tcol().Find(
		bson.M{
			"deadline": bson.M{
				"$gt": now,
			},
		}).One(&taskAssociationItem)
	if err != nil {
		return nil, err
	}

	task, err := db.getTask(taskAssociationItem)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (db *Database) getTask(t taskAssociation) (*models.Task, error) {
	var task models.Task

	// Query for the action
	err := db.acol().Find(
		bson.M{
			"_id": t.ActionID,
		}).One(&task.Action)
	if err != nil {
		return nil, err
	}

	// Query for the user
	err = db.ucol().Find(
		bson.M{
			"_id": t.UserID,
		}).One(&task.User)
	if err != nil {
		return nil, err
	}

	// Set the deadline
	task.Deadline = t.Deadline

	return &task, nil
}

func (db *Database) GetUsers() ([]models.User, error) {
	var users []models.User
	err := db.ucol().Find(nil).All(&users)
	return users, err
}

func (db *Database) GetActions() ([]models.Action, error) {
	var actions []models.Action
	err := db.acol().Find(nil).All(&actions)
	return actions, err
}

func (db *Database) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	var tasksAssoc []taskAssociation
	if err := db.tcol().Find(nil).All(&tasksAssoc); err != nil {
		return nil, err
	}
	for _, taskAssoc := range tasksAssoc {
		task, err := db.getTask(taskAssoc)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)

	}
	return tasks, nil
}

func (db *Database) AddUser(u models.User) error {
	ucol := db.ucol()
	return ucol.Insert(&u)

}

func (db *Database) DeleteUser(name string) error {
	return db.ucol().Remove(
		bson.M{
			"fullname": name,
		})
}

func (db *Database) AddAction(a models.Action) error {
	acol := db.acol()
	return acol.Insert(&a)
}

func (db *Database) DeleteAction(description string) error {

	return db.acol().Remove(
		bson.M{
			"description": description,
		})

}

func (db *Database) AddTask(u models.Task) error {
	taskAssoc := taskAssociation{
		UserID:   bson.ObjectId(u.User.ID),
		ActionID: bson.ObjectId(u.Action.ID),
		Deadline: u.Deadline,
	}

	return db.tcol().Insert(&taskAssoc)

}
