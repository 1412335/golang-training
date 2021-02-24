package handler

import (
	"context"

	// "github.com/micro/dev/model"
	// "github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	// "github.com/gosimple/slug"

	tags "tags/proto"
)

type Tags struct {
	db model.Model
}

func NewTags() *Tags {
	slugIndex := model.ByEquality("slug")
	slugIndex.Order.Type = model.OrderTypeUnordered

	t := &Tags{
		db: model.NewModel(
			model.WithKey("slug"),
			model.WithIndexes(model.ByEquality("type"), slugIndex),
		),
	}
	t.db.Register(tags.Tag{})
	return t
}

func (t *Tags) Add(ctx context.Context, req *tags.AddRequest, rsp *tags.AddResponse) error {
	logger.Info("Received tags.Add request")
	return nil
}

func (t *Tags) Remove(ctx context.Context, req *tags.RemoveRequest, rsp *tags.RemoveResponse) error {
	logger.Info("Received tags.Remove request")
	return nil
}

func (t *Tags) Update(ctx context.Context, req *tags.UpdateRequest, rsp *tags.UpdateResponse) error {
	logger.Info("Received tags.Update request")
	return nil
}

func (t *Tags) List(ctx context.Context, req *tags.ListRequest, rsp *tags.ListResponse) error {
	logger.Info("Received tags.List request")
	return nil
}
