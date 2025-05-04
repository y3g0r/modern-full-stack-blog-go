package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/repo"
)

type JamsRepository interface {
	CreateJam(context.Context, domain.Jam) error
	GetAllJams(context.Context, repo.GetAllJamsParams) ([]domain.Jam, error)
}

type Jams struct {
	repo      JamsRepository
	nextJamId int64
	lock      sync.Mutex
}

func NewJams(repo JamsRepository) *Jams {
	return &Jams{
		repo:      repo,
		nextJamId: 1,
	}
}

type CreateJamParams struct {
	CreatedByUserId           string
	Name                      string
	StartTimestamp            time.Time
	EndTimestamp              time.Time
	Location                  string
	ParticipantEmailAddresses []string
}

type CreateJamResult struct {
	JamId string
}

func (j *Jams) getUserPrimaryEmailAddress(ctx context.Context, usrId string) (string, error) {
	var address string

	// TODO: clean architecture, use interface & application types
	usr, err := user.Get(ctx, usrId)
	if err != nil {
		return address, err
	}

	for _, a := range usr.EmailAddresses {
		if a.ID == *usr.PrimaryEmailAddressID {
			address = a.EmailAddress
			break
		}
	}
	if address == "" {
		return address, fmt.Errorf("Expected every user to have PrimaryEmailAddress, looks like need to learn more how Clerk works")
	}
	return address, nil
}

func (j *Jams) CreateJam(ctx context.Context, p CreateJamParams) (CreateJamResult, error) {
	j.lock.Lock()
	defer j.lock.Unlock()

	jamId := strconv.FormatInt(j.nextJamId, 10)
	defer func() { j.nextJamId++ }()

	organizerEmail, err := j.getUserPrimaryEmailAddress(ctx, p.CreatedByUserId)
	if err != nil {
		return CreateJamResult{}, err
	}

	var organizerInParticipants bool
	participants := make([]domain.Participant, len(p.ParticipantEmailAddresses))
	for i, emailAddress := range p.ParticipantEmailAddresses {
		if emailAddress == organizerEmail {
			organizerInParticipants = true
		}
		participants[i] = domain.Participant{EmailAddress: emailAddress}
	}
	if !organizerInParticipants {
		participants = append(participants, domain.Participant{EmailAddress: organizerEmail})
	}

	jam := domain.Jam{
		ID:             jamId,
		CreatedBy:      organizerEmail,
		Name:           p.Name,
		StartTimestamp: p.StartTimestamp,
		EndTimestamp:   p.EndTimestamp,
		Location:       p.Location,
		Participants:   participants,
	}
	defer func() { j.nextJamId++ }()

	j.repo.CreateJam(ctx, jam)
	return CreateJamResult{JamId: jamId}, nil
}

type GetAllJamsParams struct {
	UserId string
}

type GetAllJamsResult struct {
	Jams []domain.Jam
}

func (j *Jams) GetAllJams(ctx context.Context, p GetAllJamsParams) (GetAllJamsResult, error) {
	// TODO: clean architecture, use interface & application types
	usr, err := user.Get(ctx, p.UserId)
	if err != nil {
		return GetAllJamsResult{}, err
	}

	jams, err := j.repo.GetAllJams(ctx, repo.GetAllJamsParams{
		UserEmailAddress: usr.EmailAddresses[0].EmailAddress, // FIXME
	})
	if err != nil {
		return GetAllJamsResult{}, err
	}

	return GetAllJamsResult{
		Jams: jams,
	}, nil
}
