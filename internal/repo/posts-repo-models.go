package repo

import "database/sql"

type PostRecord struct {
	ID      int     `db:"id"`
	Title   string  `db:"title"`
	Content sql.NullString `db:"content"`
}