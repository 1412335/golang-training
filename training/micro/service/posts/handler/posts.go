package handler

import (
	"context"

	"github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	// model "github.com/micro/micro/v3/service/model"

	"github.com/gosimple/slug"

	posts "posts/proto"
)

type Posts struct {
	db model.Model
}

func NewPosts() *Posts {
	createdIndex := model.ByEquality("created")
	createdIndex.Order.Type = model.OrderTypeDesc

	// create model
	db := model.NewModel(
		model.WithIndexes(model.ByEquality("slug"), createdIndex),
	)
	// register the post instance
	db.Register(new(posts.Post))

	return &Posts{
		db: db,
	}
}

func (p *Posts) Save(context.Context, *SaveRequest, *SaveResponse) error {
	logger.Info("Received Posts.Save request")
	return nil
}

// Index returns the posts index without content
func (p *Posts) Index(context.Context, *IndexRequest, *IndexResponse) error {
	return nil
}

// Query currently only supports read by slug or timestamp, no listing.
func (p *Posts) Query(context.Context, *QueryRequest, *QueryResponse) error {
	return nil
}

func (p *Posts) Delete(context.Context, *DeleteRequest, *DeleteResponse) error {
	return nil
}
