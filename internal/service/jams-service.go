package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type JamsRepository interface {
	CreateJam(context.Context, domain.Jam) error
	GetAllJams(context.Context) ([]domain.Jam, error)
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

func (j *Jams) CreateJam(ctx context.Context, p CreateJamParams) (CreateJamResult, error) {
	j.lock.Lock()
	defer j.lock.Unlock()

	jamId := strconv.FormatInt(j.nextJamId, 10)
	defer func() { j.nextJamId++ }()

	// TODO: clean architecture, use interface & application types
	usr, err := user.Get(ctx, p.CreatedByUserId)
	if err != nil {
		return CreateJamResult{}, err
	}

	participants := make([]domain.Participant, len(p.ParticipantEmailAddresses))
	for i, emailAddress := range p.ParticipantEmailAddresses {
		participants[i] = domain.Participant{EmailAddress: emailAddress}
	}

	jam := domain.Jam{
		ID:             jamId,
		CreatedBy:      usr.EmailAddresses[0].EmailAddress, // FIXME
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
	jams, err := j.repo.GetAllJams(ctx)
	if err != nil {
		return GetAllJamsResult{}, err
	}

	return GetAllJamsResult{
		Jams: jams,
	}, nil
}
