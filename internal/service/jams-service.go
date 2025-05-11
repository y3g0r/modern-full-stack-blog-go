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
	CreateJamInviteResponse(context.Context, repo.CreateJamInviteResponseParams) error
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
	JamId domain.JamId
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
	return CreateJamResult{JamId: domain.JamId(jamId)}, nil
}

type GetAllJamsParams struct {
	UserId string
}

type GetAllJamsResult struct {
	Jams []domain.Jam
}

func (j *Jams) GetAllJams(ctx context.Context, p GetAllJamsParams) (GetAllJamsResult, error) {
	userEmailAddress, err := j.getUserPrimaryEmailAddress(ctx, p.UserId)
	if err != nil {
		return GetAllJamsResult{}, err
	}

	jams, err := j.repo.GetAllJams(ctx, repo.GetAllJamsParams{
		UserEmailAddress: userEmailAddress,
	})
	if err != nil {
		return GetAllJamsResult{}, err
	}

	return GetAllJamsResult{
		Jams: jams,
	}, nil
}

type ResponseToJamInviteParams struct {
	JamId    domain.JamId
	UserId   string
	Response domain.InviteResponse
}

func (j *Jams) RespondToJamInvite(ctx context.Context, p ResponseToJamInviteParams) error {
	inviteeEmailAddress, err := j.getUserPrimaryEmailAddress(ctx, p.UserId)
	if err != nil {
		return err
	}

	err = j.repo.CreateJamInviteResponse(ctx, repo.CreateJamInviteResponseParams{
		JamId:                   p.JamId,
		ParticipantEmailAddress: inviteeEmailAddress,
		Response:                p.Response,
		ResponseTimestamp:       time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
