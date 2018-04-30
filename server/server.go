package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/think-it-labs/actionizer/models"

	log "github.com/sirupsen/logrus"
	"github.com/think-it-labs/actionizer/datastore"
)

type enforceRequest struct {
	Names []string  `json:"names"`
	When  time.Time `json:"when"`
}

func (s *Server) getCurrentTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks := s.DS.GetCurrentTasks()
	json.NewEncoder(w).Encode(tasks)
}

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, _ := s.DS.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func (s *Server) getActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	actions, _ := s.DS.GetActions()
	json.NewEncoder(w).Encode(actions)
}

func (s *Server) getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasks, _ := s.DS.GetTasks()
	json.NewEncoder(w).Encode(tasks)
}

func (s *Server) enforce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var payload enforceRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	log.Debugf("Enforce request: %+v", payload)

	allUsers, _ := s.DS.GetUsers()
	var users []models.User
	for _, name := range payload.Names {
		if user, ok := allUsers[name]; ok {
			users = append(users, user)
		} else {
			log.Warnf("User %q not found.", name)
		}
	}

	err = s.EnforceHandler(users, payload.When)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string
	}{
		Message: fmt.Sprintf("Enforced action to %d user", len(users)),
	}
	json.NewEncoder(w).Encode(&response)
}

type Server struct {
	Host           string
	Port           int
	DS             datastore.DataStore
	EnforceHandler func(users []models.User, when time.Time) error
}

func (s Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	http.HandleFunc("/current_tasks", s.getCurrentTask)
	http.HandleFunc("/users", s.getUsers)
	http.HandleFunc("/actions", s.getActions)
	http.HandleFunc("/tasks", s.getTasks)
	http.HandleFunc("/tasks/enforce", s.enforce)

	http.Handle("/", http.FileServer(http.Dir("public/")))

	return http.ListenAndServe(addr, nil)
}
