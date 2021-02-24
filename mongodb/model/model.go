package model

import "context"

type Model interface {
	Distinct(ctx context.Context, field string, filter interface{}) ([]interface{}, error)
	AggregateRaw(ctx context.Context, pipeline interface{}) ([]interface{}, error)
	Counts(ctx context.Context, filter interface{}) (int64, error)
}
