//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml openapi.yaml

package api

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

type BlogApi struct {
	posts     *service.Posts
	jams      []Jam
	nextJamId int64
	lock      sync.Mutex
}

var _ StrictServerInterface = (*BlogApi)(nil)

var duration time.Duration

func init() {
	var err error
	duration, err = time.ParseDuration("2h")
	if err != nil {
		panic(err)
	}
}

func NewBlog(posts *service.Posts) *BlogApi {
	return &BlogApi{
		posts: posts,
		jams: []Jam{
			{
				Id:                    1,
				CreatedBy:             "dummy",
				Name:                  "Hardcoded in API",
				StartTimestampSeconds: time.Now().Unix(),
				EndTimestampSeconds:   time.Now().Add(duration).Unix(),
			},
		},
		nextJamId: 2,
	}
}

// GetJams implements StrictServerInterface.
func (b *BlogApi) GetJams(ctx context.Context, request GetJamsRequestObject) (GetJamsResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return GetJams401Response{}, nil
	}

	result := []Jam{}
	for _, jam := range b.jams {
		result = append(result, Jam{
			Id:                    jam.Id,
			CreatedBy:             claims.Subject,
			Name:                  jam.Name,
			StartTimestampSeconds: jam.StartTimestampSeconds,
			EndTimestampSeconds:   jam.EndTimestampSeconds,
		})
	}

	return GetJams200JSONResponse(result), nil
}

// CreateJam implements StrictServerInterface.
func (b *BlogApi) CreateJam(ctx context.Context, request CreateJamRequestObject) (CreateJamResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return CreateJam401Response{}, nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	newJam := Jam{
		Id:                    b.nextJamId,
		CreatedBy:             claims.Subject,
		Name:                  request.Body.Name,
		StartTimestampSeconds: request.Body.StartTimestampSeconds,
		EndTimestampSeconds:   request.Body.EndTimestampSeconds,
		Location:              request.Body.Location,
		Participants:          request.Body.Participants,
	}
	b.nextJamId += 1

	b.jams = append(b.jams, newJam)
	return CreateJam201JSONResponse{Id: newJam.Id}, nil
}

// CreatePost implements StrictServerInterface.
func (b *BlogApi) CreatePost(ctx context.Context, request CreatePostRequestObject) (CreatePostResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return CreatePost201JSONResponse{}, fmt.Errorf("missing authentication claims in CreatePost request context, is authentication middleware misconfigured?")
	}
	usr, err := user.Get(ctx, claims.Subject)
	if err != nil {
		return CreatePost201JSONResponse{}, err
	}

	slog.Info(fmt.Sprintf(`{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned))
	slog.Info(fmt.Sprintf("%#v", usr))

	post, err := b.posts.CreatePost(service.CreatePostParams{
		CreatedBy: domain.UserId(claims.Subject),
		Title:     request.Body.Title,
		Content:   request.Body.Content,
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
			Author:  (*string)(&post.CreatedBy),
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
