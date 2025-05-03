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

// GetAllJams implements service.JamsRepository.
func (j *Jams) GetAllJams(context.Context) ([]domain.Jam, error) {
	return j.jams, nil
}
