package service

import (
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
)

type PostsRepo interface {
	CreatePost(post domain.Post) error
	GetPost(id int) (domain.Post, error)
	UpdatePost(postId int, params UpdatePostParams) error
	DeletePost(id int) error
	GetPosts() ([]domain.Post, error)
}

type Posts struct {
	postsRepo PostsRepo
	lastId  int
}

func NewPostsService(blogDao PostsRepo) *Posts {
	return &Posts{
		postsRepo: blogDao,
	}
}

func (s *Posts) CreatePost(params CreatePostParams) (domain.Post, error) {
	s.lastId++
	post := domain.Post{
		ID:        s.lastId,
		CreatedBy: params.CreatedBy,
		Title:     params.Title,
		Content:   params.Content,
	}
	if err := s.postsRepo.CreatePost(post); err != nil {
		return domain.Post{}, err
	}
	return post, nil
}

func (s *Posts) GetPost(id int) (domain.Post, error) {
	post, err := s.postsRepo.GetPost(id)
	if err != nil {
		return domain.Post{}, err
	}
	return post, nil
}

func (s *Posts) UpdatePost(postId int, params UpdatePostParams) error {
	if err := s.postsRepo.UpdatePost(postId, params); err != nil {
		return err
	}
	return nil
}

func (s *Posts) DeletePost(id int) error {
	if err := s.postsRepo.DeletePost(id); err != nil {
		return err
	}
	return nil
}

func (s *Posts) GetPosts() ([]domain.Post, error) {
	posts, err := s.postsRepo.GetPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
