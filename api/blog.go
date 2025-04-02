//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml openapi.yaml

package api

import (
	"context"

	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

type BlogApi struct {
	posts *service.Posts
}

var _ StrictServerInterface = (*BlogApi)(nil)

func NewBlog(posts *service.Posts) *BlogApi {
	return &BlogApi{
		posts: posts,
	}
}

// CreatePost implements StrictServerInterface.
func (b *BlogApi) CreatePost(ctx context.Context, request CreatePostRequestObject) (CreatePostResponseObject, error) {
	post, err := b.posts.CreatePost(service.CreatePostParams{
		Title:   request.Body.Title,
		Content: request.Body.Content,
	})
	if err != nil {
		return CreatePost201JSONResponse{}, err
	}

	return CreatePost201JSONResponse{Id: post.ID}, nil
}

// DeletePost implements StrictServerInterface.
func (b *BlogApi) DeletePost(ctx context.Context, request DeletePostRequestObject) (DeletePostResponseObject, error) {
	err := b.posts.DeletePost(request.Id)

	return DeletePost204Response{}, err
}

// GetPost implements StrictServerInterface.
func (b *BlogApi) GetPost(ctx context.Context, request GetPostRequestObject) (GetPostResponseObject, error) {
	panic("unimplemented")
}

// GetPosts implements StrictServerInterface.
func (b *BlogApi) GetPosts(ctx context.Context, request GetPostsRequestObject) (GetPostsResponseObject, error) {
	postList := make([]Post, 0)
	posts, err := b.posts.GetPosts()
	if err != nil {
		return GetPosts200JSONResponse(postList), err
	}

	for _, post := range posts {
		postList = append(postList, Post{
			Id:      &post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}

	return GetPosts200JSONResponse(postList), nil
}

// UpdatePost implements StrictServerInterface.
func (b *BlogApi) UpdatePost(ctx context.Context, request UpdatePostRequestObject) (UpdatePostResponseObject, error) {
	panic("unimplemented")
}
