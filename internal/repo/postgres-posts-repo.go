package repo

import (
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

type PostgresPostsRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewPostgresPostsRepo(db *sqlx.DB, logger *slog.Logger) *PostgresPostsRepo {
	return &PostgresPostsRepo{
		db:     db,
		logger: logger,
	}
}

// TODO: probably will need to add context parameter to all methods
// TODO: will need to introduce explicit transaction management
// CreatePost implements service.PostsRepo.
func (r *PostgresPostsRepo) CreatePost(post domain.Post) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	var content sql.NullString
	if post.Content != nil {
		content = sql.NullString{String: *post.Content, Valid: true}
	} else {
		content = sql.NullString{Valid: false}
	}

	_, err = tx.Exec("INSERT INTO posts (id, created_by, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.CreatedBy, post.Title, content)
	if err != nil {
		r.logger.Error("Error on attempt to insert record into DB: " + err.Error())
		return err
	}
	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error on commiting tx to insert new post: " + err.Error())
		return err
	}
	return nil
}

// DeletePost implements service.PostsRepo.
func (r *PostgresPostsRepo) DeletePost(id int) error {
	panic("unimplemented")
}

// GetPost implements service.PostsRepo.
func (r *PostgresPostsRepo) GetPost(id int) (domain.Post, error) {
	panic("unimplemented")
}

// GetPosts implements service.PostsRepo.
func (r *PostgresPostsRepo) GetPosts() ([]domain.Post, error) {
	rows, err := r.db.Queryx("SELECT id, created_by, title, content FROM posts")
	if err != nil {
		return []domain.Post{}, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0)
	for rows.Next() {
		var record PostRecord
		err = rows.StructScan(&record)
		if err != nil {
			r.logger.Error("Failed to struct scan: " + err.Error())
			return []domain.Post{}, err
		}

		var content *string
		if record.Content.Valid {
			content = &record.Content.String
		}
		posts = append(posts, domain.Post{
			ID:        record.ID,
			CreatedBy: domain.UserId(record.CreatedBy),
			Title:     record.Title,
			Content:   content,
		})
	}
	return posts, nil

}

// UpdatePost implements service.PostsRepo.
func (r *PostgresPostsRepo) UpdatePost(postId int, params service.UpdatePostParams) error {
	panic("unimplemented")
}
