package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/logger"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

type PostgresPostsRepo struct {
	db     *sqlx.DB
	logger logger.Interface
}

func NewPostgresPostsRepo(db *sqlx.DB, logger logger.Interface) *PostgresPostsRepo {
	return &PostgresPostsRepo{
		db: db,
	}
}

// TODO: probably will need to add context parameter to all methods
// TODO: will need to introduce explicit transaction management
// CreatePost implements service.PostsRepo.
func (r *PostgresPostsRepo) CreatePost(post domain.Post) error {
	panic("unimplemented")
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
	rows, err := r.db.Queryx("SELECT id, title, content FROM posts")
	if err != nil {
		return []domain.Post{}, err
	}
	defer rows.Close()

	posts := make([]domain.Post, 0)
	for rows.Next() {
		var record PostRecord
		err = rows.StructScan(&record)
		if err != nil {
			r.logger.Error("Failed to struct scan " + err.Error())
			return []domain.Post{}, err
		}

		var content *string
		if record.Content.Valid {
			content = &record.Content.String
		}
		posts = append(posts, domain.Post{
			ID:      record.ID,
			Title:   record.Title,
			Content: content,
		})
	}
	return posts, nil

}

// UpdatePost implements service.PostsRepo.
func (r *PostgresPostsRepo) UpdatePost(postId int, params service.UpdatePostParams) error {
	panic("unimplemented")
}
