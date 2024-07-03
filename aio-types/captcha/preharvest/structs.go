package preharvest

import (
	"regexp"
	"time"

	"github.com/Copped-Inc/aio-types/modules"
	"github.com/infinitare/disgo"
)

type (
	Task struct {
		ID      string          `json:"id" firestore:"-"`
		User_ID disgo.Snowflake `json:"user_id" firestore:"user_id"`
		State   State           `json:"state" firestore:"state"`
		Task_Create
		Uses int `json:"uses,omitempty" firestore:"uses,omitempty"`
	}

	State int

	Task_Create struct {
		Routine  bool         `json:"routine,omitempty" firestore:"routine,omitempty"`
		Date     time.Time    `json:"date,omitempty" firestore:"date,omitempty"`
		Site     modules.Site `json:"site" firestore:"site"`
		Schedule string       `json:"schedule,omitempty" firestore:"schedule,omitempty"`
	}

	Task_Edit struct {
		State    State     `json:"state,omitempty"`
		Uses     int       `json:"uses,omitempty"`
		Schedule string    `json:"schedule,omitempty"`
		Date     time.Time `json:"date,omitempty"`
	}
)

const (
	Running State = iota + 1
	Stopped
)

var Schedule_Pattern = regexp.MustCompile(`^(?si:(?:every (?:(2|3|4|5|6) (day|week|month)(?:s)?|(day|week|month)))|(daily|weekly|monthly))$`)
