package repo

import (
	"context"
	"sync"

	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type Jams struct {
	jams []domain.Jam
	lock sync.Mutex
}

func NewInMemoryJams() *Jams {
	return &Jams{
		// TODO: check if there is a difference between leaving it nil, new, make or literal
		jams: []domain.Jam{},
	}
}

// CreateJam implements service.JamsRepository.
func (j *Jams) CreateJam(ctx context.Context, jam domain.Jam) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	j.jams = append(j.jams, jam)
	return nil
}

type GetAllJamsParams struct {
	UserEmailAddress string
}

// GetAllJams implements service.JamsRepository.
func (j *Jams) GetAllJams(ctx context.Context, p GetAllJamsParams) ([]domain.Jam, error) {
	result := []domain.Jam{}
	for _, jam := range j.jams {
		if jam.CreatedBy == p.UserEmailAddress {
			result = append(result, jam)
			continue
		}

		for _, prt := range jam.Participants {
			if prt.EmailAddress == p.UserEmailAddress {
				result = append(result, jam)
				break
			}
		}
	}
	return result, nil
}
