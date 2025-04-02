package service

import "database/sql"

type CreatePostParams struct {
	Title   string
	Content *string
}

type UpdatePostParams struct {
	Title   sql.NullString
	Content sql.Null[sql.NullString]
}
