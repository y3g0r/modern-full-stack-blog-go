package repo

import (
	"fmt"

	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
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
	dao.posts[post.ID] = PostRecord(post)
	return nil
}

func (dao *InMemoryPostsRepo) GetPost(id int) (domain.Post, error) {
	post, exists := dao.posts[id]
	if !exists {
		return domain.Post{}, fmt.Errorf("post with id %d not found", id)
	}
	return domain.Post(post), nil
}

func (dao *InMemoryPostsRepo) UpdatePost(postId int, params service.UpdatePostParams) error {
	var post PostRecord
	var exists bool
	if post, exists = dao.posts[postId]; !exists {
		return fmt.Errorf("post with id %d not found", postId)
	}

	if params.Title.Valid {
		post.Title = params.Title.String
	}

	if params.Content.Valid {
		if params.Content.V.Valid {
			post.Content = &params.Content.V.String
		} else {
			post.Content = nil
		}
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
		posts = append(posts, domain.Post(post))
	}
	return posts, nil
}
