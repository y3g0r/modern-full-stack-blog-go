//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=cfg.yaml openapi.yaml

package api

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/domain"
	"github.com/y3g0r/modern-full-stack-blog-go/internal/service"
)

type Api struct {
	posts       *service.Posts
	jamsService *service.Jams
}

var _ StrictServerInterface = (*Api)(nil)

func NewApi(posts *service.Posts, jams *service.Jams) *Api {
	return &Api{
		posts:       posts,
		jamsService: jams,
	}
}

// GetJams implements StrictServerInterface.
func (b *Api) GetJams(ctx context.Context, request GetJamsRequestObject) (GetJamsResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return GetJams401Response{}, nil
	}

	result, err := b.jamsService.GetAllJams(ctx, service.GetAllJamsParams{
		UserId: claims.Subject,
	})
	if err != nil {
		return GetJams200JSONResponse{}, err
	}

	response := make([]Jam, len(result.Jams))
	for i, jam := range result.Jams {
		id, err := strconv.ParseInt(jam.ID, 10, 64)
		if err != nil {
			return GetJams200JSONResponse{}, err
		}
		participants := make([]Participant, len(jam.Participants))
		for i, p := range jam.Participants {
			accepted := Accepted
			declined := Declined
			var response *ParticipantResponse
			if r := rand.Float32(); r < 0.33 {
				response = &accepted
			} else if r < 0.66 {
				response = &declined
			}
			// else leave it nil
			participants[i] = Participant{
				Email:    p.EmailAddress,
				Response: response,
			}
		}

		response[i] = Jam{
			Id:                    id,
			CreatedBy:             jam.CreatedBy,
			Name:                  jam.Name,
			StartTimestampSeconds: jam.StartTimestamp.Unix(),
			EndTimestampSeconds:   jam.EndTimestamp.Unix(),
			Location:              jam.Location,
			Participants:          participants,
		}
	}

	return GetJams200JSONResponse(response), nil
}

// CreateJam implements StrictServerInterface.
func (b *Api) CreateJam(ctx context.Context, request CreateJamRequestObject) (CreateJamResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return CreateJam401Response{}, nil
	}

	participantsEmails := make([]string, len(request.Body.Participants))
	for i, p := range request.Body.Participants {
		participantsEmails[i] = p.Email
	}

	result, err := b.jamsService.CreateJam(ctx, service.CreateJamParams{
		CreatedByUserId:           claims.Subject,
		Name:                      request.Body.Name,
		StartTimestamp:            time.Unix(request.Body.StartTimestampSeconds, 0),
		EndTimestamp:              time.Unix(request.Body.EndTimestampSeconds, 0),
		Location:                  request.Body.Location,
		ParticipantEmailAddresses: participantsEmails,
	})
	if err != nil {
		return CreateJam201JSONResponse{}, err
	}
	id, err := strconv.ParseInt(result.JamId, 10, 64)
	if err != nil {
		return CreateJam201JSONResponse{}, err
	}
	return CreateJam201JSONResponse{Id: id}, err
}

// CreatePost implements StrictServerInterface.
func (b *Api) CreatePost(ctx context.Context, request CreatePostRequestObject) (CreatePostResponseObject, error) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok {
		return CreatePost201JSONResponse{}, fmt.Errorf("missing authentication claims in CreatePost request context, is authentication middleware misconfigured?")
	}
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
func (b *Api) DeletePost(ctx context.Context, request DeletePostRequestObject) (DeletePostResponseObject, error) {
	err := b.posts.DeletePost(request.Id)

	return DeletePost204Response{}, err
}

// GetPost implements StrictServerInterface.
func (b *Api) GetPost(ctx context.Context, request GetPostRequestObject) (GetPostResponseObject, error) {
	panic("unimplemented")
}

// GetPosts implements StrictServerInterface.
func (b *Api) GetPosts(ctx context.Context, request GetPostsRequestObject) (GetPostsResponseObject, error) {
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
func (b *Api) UpdatePost(ctx context.Context, request UpdatePostRequestObject) (UpdatePostResponseObject, error) {
	panic("unimplemented")
}
