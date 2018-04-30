package models

import (
	"time"

	"github.com/think-it-labs/actionizer/utils"
)

type User struct {
	Fullname string `json:"fullname" bson:"fullname,omitempty"`
	ImageURL string `json:"image_url" bson:"image_url,omitempty"`
	Remote   bool   `json:"remote"`
}

type Action struct {
	ID          string         `json:"action"        bson:"_id,omitempty"`
	Description string         `json:"message" bson:"description,omitempty"`
	Remotee     bool           `json:"remotee"`
	Enforce     bool           `json:"enforce"`
	Duration    utils.Duration `json:"duration"`
}

type Task struct {
	ID        string    `json:"id"`
	User      User      `json:"user"`
	Action    Action    `json:"action"`
	StartDate time.Time `json:"start_date"`
	Deadline  time.Time `json:"deadline"`
	Enforced  bool      `json:"enforced"`
	Done      bool      `json:"done"`
}
