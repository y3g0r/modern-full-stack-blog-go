package domain

import "time"

type Participant struct {
	EmailAddress string
	JamInviteResponse *InviteResponse
}

type JamId string

type InviteResponse string

const (
	InviteAccepted InviteResponse = "accepted"
	InviteDeclined InviteResponse = "declined"
)

type Jam struct {
	ID             string
	CreatedBy      string
	Name           string
	StartTimestamp time.Time
	EndTimestamp   time.Time
	Location       string
	Participants   []Participant
}
