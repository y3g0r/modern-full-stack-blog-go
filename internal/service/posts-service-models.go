package service

import (
	"database/sql"

	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type CreatePostParams struct {
	CreatedBy domain.UserId
	Title     string
	Content   *string
}

type UpdatePostParams struct {
	Title   sql.NullString
	Content sql.Null[sql.NullString]
}
