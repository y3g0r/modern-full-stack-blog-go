package repo

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type PostgresJamsRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewPostgresJamsRepo(db *sqlx.DB, logger *slog.Logger) *PostgresJamsRepo {
	return &PostgresJamsRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresJamsRepo) CreateJam(ctx context.Context, jam domain.Jam) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	var jamId int
	err = tx.Get(&jamId, "INSERT INTO jams (created_by, name, start_timestamp, end_timestamp, location) VALUES ($1, $2, $3, $4, $5) RETURNING id", jam.CreatedBy, jam.Name, jam.StartTimestamp, jam.EndTimestamp, jam.Location)
	if err != nil {
		r.logger.Error("Error on attempt to insert record into DB: " + err.Error())
		return err
	}
	for _, part := range jam.Participants {
		_, err := tx.Exec("INSERT INTO jam_participants (email, jam_id) VALUES ($1, $2)", part.EmailAddress, jamId)
		if err != nil {
			r.logger.Error("Error on attempt to insert record into DB: " + err.Error())
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error on commiting tx to insert new post: " + err.Error())
		return err
	}
	return nil
}

type JamRecord struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	CreatedBy      string    `db:"created_by"`
	StartTimestamp time.Time `db:"start_timestamp"`
	EndTimestamp   time.Time `db:"end_timestamp"`
	Location       string    `db:"location"`
}

func (r *PostgresJamsRepo) getJamIDsByEmail(ctx context.Context, address string) ([]int, error) {
	var jamIds []int
	err := r.db.Select(
		&jamIds, `
		SELECT j.id
		FROM jams j
		JOIN jam_participants p ON p.jam_id = j.id
		WHERE p.email = $1`,
		address,
	)
	if err != nil {
		r.logger.Error("Error running Select jam IDs query: " + err.Error())
		return []int{}, err
	}
	r.logger.Debug("Fetched jam IDs in which user is a participant: " + fmt.Sprint(jamIds))
	return jamIds, nil
}

func (r *PostgresJamsRepo) getJamRecordsByIds(ctx context.Context, ids []int) ([]JamRecord, error) {
	query, args, err := sqlx.In(`
	SELECT j.id, j.created_by, j.name, j.start_timestamp, j.end_timestamp, j.location
	FROM jams j
	WHERE j.id IN (?);`, ids)
	if err != nil {
		r.logger.Error("Error preping select jams by ids query: " + err.Error())
		return []JamRecord{}, err
	}

	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = r.db.Rebind(query)

	var jamRecords []JamRecord
	err = r.db.Select(&jamRecords, query, args...)
	if err != nil {
		r.logger.Error("Error running Select jams query: " + err.Error())
		return []JamRecord{}, err
	}
	r.logger.Debug("Fetched jams by IDs: " + fmt.Sprint(jamRecords))
	return jamRecords, nil
}

type ParticipantRecord struct {
	ID           int64  `db:"id"`
	EmailAddress string `db:"email"`
	JamId        int64  `db:"jam_id"`
}

func (r *PostgresJamsRepo) getParticipantsByJamIDs(ctx context.Context, ids []int) ([]ParticipantRecord, error) {
	query, args, err := sqlx.In(`
	SELECT p.id, p.email, p.jam_id
	FROM jam_participants p
	WHERE p.jam_id IN (?);`, ids)
	if err != nil {
		r.logger.Error("Error preping select participants by jam ids query: " + err.Error())
		return []ParticipantRecord{}, err
	}

	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = r.db.Rebind(query)

	var participants []ParticipantRecord
	err = r.db.Select(&participants, query, args...)
	if err != nil {
		r.logger.Error("Error running select participants by jam ids query: " + err.Error())
		return []ParticipantRecord{}, err
	}
	r.logger.Info("Fetched participants by Jam IDs: " + fmt.Sprint(participants))
	return participants, nil
}

func (r *PostgresJamsRepo) getAllJamsMultipleQueries(ctx context.Context, p GetAllJamsParams) ([]domain.Jam, error) {
	jamIds, err := r.getJamIDsByEmail(ctx, p.UserEmailAddress)
	if err != nil {
		return []domain.Jam{}, err
	}

	if len(jamIds) == 0 {
		return []domain.Jam{}, nil
	}

	jamRecords, err := r.getJamRecordsByIds(ctx, jamIds)
	if err != nil {
		return []domain.Jam{}, err
	}

	participantRecords, err := r.getParticipantsByJamIDs(ctx, jamIds)
	if err != nil {
		return []domain.Jam{}, err
	}

	jams := make([]domain.Jam, len(jamRecords))
	for i, r := range jamRecords {
		var participants []domain.Participant
		for _, p := range participantRecords {
			if p.JamId == r.ID {
				participants = append(participants, domain.Participant{
					EmailAddress: p.EmailAddress,
				})
			}
		}

		jams[i] = domain.Jam{
			ID:             strconv.FormatInt(r.ID, 10),
			CreatedBy:      r.CreatedBy,
			Name:           r.Name,
			StartTimestamp: r.StartTimestamp,
			EndTimestamp:   r.EndTimestamp,
			Location:       r.Location,
			Participants:   participants,
		}
	}

	return jams, nil
}

func (r *PostgresJamsRepo) GetAllJams(ctx context.Context, p GetAllJamsParams) ([]domain.Jam, error) {

	return r.getAllJamsMultipleQueries(ctx, p)
}
