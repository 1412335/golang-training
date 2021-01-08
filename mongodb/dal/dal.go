package dal

import "context"

type DataAccessLayer interface {
	GetClient(ctx context.Context) (interface{}, error)
	Disconnect(ctx context.Context)

	Insert(ctx context.Context, collectionName string, data interface{}) (interface{}, error)
	InsertMany(ctx context.Context, collectionName string, data []interface{}) ([]interface{}, error)

	Update(ctx context.Context, collectionName string, filter, data interface{}, upsert bool) error
	UpdateMany(ctx context.Context, collectionName string, filter, data interface{}, upsert bool) error

	Delete(ctx context.Context, collectionName string, filter interface{}) error
	DeleteMany(ctx context.Context, collectionName string, filter interface{}) error

	Find(ctx context.Context, collectionName string, filter, projection, sort interface{}, offset, limit int64) ([]interface{}, error)
	FindAll(ctx context.Context, collectionName string, projection, sort interface{}) ([]interface{}, error)
	FindOne(ctx context.Context, collectionName string, filter, projection interface{}) (interface{}, error)
	Distinct(ctx context.Context, collectionName string, field string, filter interface{}) ([]interface{}, error)
	FindByID(ctx context.Context, collectionName string, id interface{}, projection interface{}) (interface{}, error)

	Aggregate(ctx context.Context, collectionName string, filter, projection, sort interface{}, offset, limit int64) ([]interface{}, error)
	AggregateCommon(ctx context.Context, collectionName string, pipeline interface{}) ([]interface{}, error)

	Counts(ctx context.Context, collectionName string, filter interface{}) (int64, error)
}
