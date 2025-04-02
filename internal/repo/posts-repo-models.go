package repo

type PostRecord struct {
	ID      int     `db:"id"`
	Title   string  `db:"title"`
	Content *string `db:"content"`
}