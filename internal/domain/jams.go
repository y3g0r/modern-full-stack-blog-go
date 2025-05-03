package domain

import "time"

type Participant struct {
	EmailAddress string
}

type Jam struct {
	ID             string
	CreatedBy      string
	Name           string
	StartTimestamp time.Time
	EndTimestamp   time.Time
	Location       string
	Participants   []Participant
}
