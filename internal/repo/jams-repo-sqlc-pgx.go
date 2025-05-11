package repo

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/repo/postgres"
)

type repo struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewJamsRepoSqlc(db *pgxpool.Pool, logger *slog.Logger) *repo {
	return &repo{
		db:     db,
		logger: logger,
	}
}

// CreateJam implements service.JamsRepository.
func (r *repo) CreateJam(ctx context.Context, jam domain.Jam) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("Error on attempt to create transaction: " + err.Error())
	}
	defer tx.Rollback(ctx)

	queries := postgres.New(tx)
	created, err := queries.CreateJam(ctx, postgres.CreateJamParams{
		CreatedBy:      jam.CreatedBy,
		Name:           pgtype.Text{String: jam.Name, Valid: true},
		StartTimestamp: pgtype.Timestamp{Time: jam.StartTimestamp, Valid: true},
		EndTimestamp:   pgtype.Timestamp{Time: jam.EndTimestamp, Valid: true},
		Location:       pgtype.Text{String: jam.Location, Valid: true},
	})
	if err != nil {
		r.logger.Error("Error on attempt to insert jam DB: " + err.Error())
		return err
	}

	for _, p := range jam.Participants {
		_, err := queries.CreateJamParticipant(ctx, postgres.CreateJamParticipantParams{
			Email: p.EmailAddress,
			JamID: created.ID,
		})
		if err != nil {
			r.logger.Error("Error on attempt to insert jam participant into DB: " + err.Error())
			return err
		}
	}

	return tx.Commit(ctx)
}

// GetAllJams implements service.JamsRepository.
func (r *repo) GetAllJams(ctx context.Context, p GetAllJamsParams) ([]domain.Jam, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error acquiring connection from the pool: " + err.Error())
		return []domain.Jam{}, err
	}
	queries := postgres.New(conn)
	jamIds, err := queries.GetJamIdsByParticipantEmail(ctx, p.UserEmailAddress)
	if err != nil {
		return []domain.Jam{}, err
	}

	if len(jamIds) == 0 {
		return []domain.Jam{}, nil
	}

	jamRecords, err := queries.GetJamsByIDs(ctx, jamIds)
	if err != nil {
		return []domain.Jam{}, err
	}

	participantRecords, err := queries.GetParticipantsByJamIDs(ctx, jamIds)
	if err != nil {
		return []domain.Jam{}, err
	}

	jams := make([]domain.Jam, len(jamRecords))
	for i, r := range jamRecords {
		var participants []domain.Participant
		for _, p := range participantRecords {
			if p.JamID == r.ID {
				participants = append(participants, domain.Participant{
					EmailAddress: p.Email,
				})
			}
		}

		jams[i] = domain.Jam{
			ID:             strconv.FormatInt(int64(r.ID), 10),
			CreatedBy:      r.CreatedBy,
			Name:           r.Name.String,
			StartTimestamp: r.StartTimestamp.Time,
			EndTimestamp:   r.EndTimestamp.Time,
			Location:       r.Location.String,
			Participants:   participants,
		}
	}

	return jams, nil
}
