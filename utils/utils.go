package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type Duration time.Duration

var unitsMultipliers = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
	"d": time.Hour * 24,
	"w": time.Hour * 24 * 7,
}

func (this *Duration) UnmarshalJSON(b []byte) error {
	match, _ := regexp.MatchString(`^['"][0-9]+[smhdw]['"]$`, string(b))
	if !match {
		return errors.New("Cannot unmarshal value for type utils.Duration")
	}
	v, _ := strconv.Unquote(string(b))
	unit := strings.ToLower(string(v[len(v)-1]))
	val, _ := strconv.Atoi(v[:len(v)-1])
	multiplier := unitsMultipliers[unit]

	*this = Duration(val) * Duration(multiplier)
	return nil
}

func (this Duration) MarshalJSON() ([]byte, error) {
	weeks := int(this) / int(unitsMultipliers["w"])
	s := fmt.Sprintf("%dw", weeks)
	return json.Marshal(s)
}

func NextWeekStart(now time.Time) time.Time {
	weekStart := now.AddDate(0, 0, 7-int(now.Weekday()))
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
}

func StartOfWeek(now time.Time) time.Time {
	weekStart := now.AddDate(0, 0, 1-int(now.Weekday()))
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())
}

func AddFileWatcher(filename string, callback func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					callback()
				}
			case err := <-watcher.Errors:
				log.Errorf("error:", err)
			}
		}
	}()

	err = watcher.Add(filename)
	if err != nil {
		return fmt.Errorf("Couldn't add watcher for file %q", filename)
	}

	return nil
}
