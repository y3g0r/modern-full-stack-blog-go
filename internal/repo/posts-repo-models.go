package repo

import "database/sql"

type PostRecord struct {
	ID        int            `db:"id"`
	CreatedBy string         `db:"created_by"`
	Title     string         `db:"title"`
	Content   sql.NullString `db:"content"`
}
