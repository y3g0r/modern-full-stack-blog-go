package domain

type UserId string

type Post struct {
	ID        int
	CreatedBy UserId
	Title     string
	Content   *string
}
