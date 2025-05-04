package repo

import "database/sql"

type PostRecord struct {
	ID        int            `db:"id"`
	CreatedBy string         `db:"created_by"`
	Title     string         `db:"title"`
	Content   sql.NullString `db:"content"`
}

type UpdatePostParams struct {
	Title   sql.NullString
	// Nesting is intentional here.
	// If top level is not null it means update was requested.
	// If value is null it means user wants to remove the content.
	// If top-level is null it means no update to the field intended.
	Content sql.Null[sql.NullString]
}

