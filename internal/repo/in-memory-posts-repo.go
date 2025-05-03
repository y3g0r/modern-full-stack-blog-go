package repo

import (
	"database/sql"
	"fmt"

	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type InMemoryPostsRepo struct {
	posts map[int]PostRecord
}

func NewInMemoryPostsRepo() *InMemoryPostsRepo {
	return &InMemoryPostsRepo{
		posts: make(map[int]PostRecord),
	}
}

func (dao *InMemoryPostsRepo) CreatePost(post domain.Post) error {
	var content sql.NullString
	if post.Content != nil {
		content = sql.NullString{String: *post.Content, Valid: true}
	} else {
		content = sql.NullString{Valid: false}
	}

	dao.posts[post.ID] = PostRecord{
		ID:      post.ID,
		Title:   post.Title,
		Content: content,
	}
	return nil
}

func (dao *InMemoryPostsRepo) GetPost(id int) (domain.Post, error) {
	post, exists := dao.posts[id]
	if !exists {
		return domain.Post{}, fmt.Errorf("post with id %d not found", id)
	}

	var content *string
	if post.Content.Valid {
		content = &post.Content.String
	}

	return domain.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: content,
	}, nil
}

func (dao *InMemoryPostsRepo) UpdatePost(postId int, params UpdatePostParams) error {
	var post PostRecord
	var exists bool
	if post, exists = dao.posts[postId]; !exists {
		return fmt.Errorf("post with id %d not found", postId)
	}

	if params.Title.Valid {
		post.Title = params.Title.String
	}

	if params.Content.Valid {
		post.Content = params.Content.V
	}
	dao.posts[postId] = post
	return nil
}

func (dao *InMemoryPostsRepo) DeletePost(id int) error {
	_, exists := dao.posts[id]
	if !exists {
		return fmt.Errorf("post with id %d not found", id)
	}
	delete(dao.posts, id)
	return nil
}

func (dao *InMemoryPostsRepo) GetPosts() ([]domain.Post, error) {
	posts := make([]domain.Post, 0, len(dao.posts))
	for _, post := range dao.posts {
		var content *string
		if post.Content.Valid {
			content = &post.Content.String
		}
		posts = append(posts, domain.Post{
			ID:      post.ID,
			Title:   post.Title,
			Content: content,
		})
	}
	return posts, nil
}
