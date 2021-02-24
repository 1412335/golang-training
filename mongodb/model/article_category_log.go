package model

import (
	"context"
	"encoding/json"
	"golang-training/mongodb/dal"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ArticleCategory struct của project
type ArticleCategoryLog struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	SourceID  string             `bson:"source_id" json:"source_id"`
	Data      string             `bson:"data" json:"data"`
	Partition int32              `bson:"partition" json:"partition"`
	Offset    int64              `bson:"offset" json:"offset"`
	Status    int16              `bson:"status" json:"status"`
	StatusLog string             `bson:"status_log" json:"status_log"`
	DstID     string             `bson:"dst_id" json:"dst_id"`
	Created   string             `bson:"created" json:"created"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Modified  string             `bson:"modified" json:"modified"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ArticleCategoryLogCollection struct {
	dal        dal.DataAccessLayer
	collection string
}

func NewArticleCategoryLogCollection(dal dal.DataAccessLayer) (*ArticleCategoryLogCollection, error) {
	return &ArticleCategoryLogCollection{
		dal:        dal,
		collection: "article_category_log",
	}, nil
}

// dataConvert xử lý data về đúng struct
func (o *ArticleCategoryLogCollection) dataConvert(listDocument []interface{}) (data []ArticleCategory, err error) {
	byteData, err := json.Marshal(listDocument)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(byteData, &data); err != nil {
		return nil, err
	}
	return data, nil
}

//ObjectConvert xử lý data về đúng struct
func (o *ArticleCategoryLogCollection) objectConvert(obj interface{}) (*ArticleCategory, error) {
	byteData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	data := &ArticleCategory{}
	if err := json.Unmarshal(byteData, data); err != nil {
		// return nil, err
	}
	return data, nil
}

// Create a document in the collection.
func (o *ArticleCategoryLogCollection) Insert(ctx context.Context, data *ArticleCategory) (*ArticleCategory, error) {
	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	_, err := o.dal.Insert(ctx, o.collection, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Create many documents in the collection.
func (o *ArticleCategoryLogCollection) InsertMany(ctx context.Context, data []ArticleCategory) ([]ArticleCategory, error) {
	items := make([]interface{}, len(data))
	for i, item := range data {
		item.ID = primitive.NewObjectID()
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		items[i] = item
		data[i] = item
	}
	_, err := o.dal.InsertMany(ctx, o.collection, items)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Update a document in the collection.
func (o *ArticleCategoryLogCollection) Update(ctx context.Context, filter, data interface{}, upsert bool) error {
	// data.UpdatedAt = time.Now()
	err := o.dal.Update(ctx, o.collection, filter, data, upsert)
	if err != nil {
		return err
	}
	return nil
}

// Update a document in the collection.
func (o *ArticleCategoryLogCollection) UpdateMany(ctx context.Context, filter, data interface{}, upsert bool) error {
	// data.UpdatedAt = time.Now()
	err := o.dal.UpdateMany(ctx, o.collection, filter, data, upsert)
	if err != nil {
		return err
	}
	return nil
}

// Find documents in the collection.
func (o *ArticleCategoryLogCollection) Find(ctx context.Context, filter, projection, sort interface{}, offset, limit int64) ([]ArticleCategory, error) {
	data, err := o.dal.Find(ctx, o.collection, filter, projection, sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return o.dataConvert(data)
}

// Find a document by ObjectID in the collection.
func (o *ArticleCategoryLogCollection) FindByID(ctx context.Context, id interface{}, projection interface{}) (*ArticleCategory, error) {
	data, err := o.dal.FindByID(ctx, o.collection, id, projection)
	if err != nil {
		return nil, err
	}
	return o.objectConvert(data)
}

// Find a document in the collection.
func (o *ArticleCategoryLogCollection) FindOne(ctx context.Context, filter, projection interface{}) (*ArticleCategory, error) {
	data, err := o.dal.FindOne(ctx, o.collection, filter, projection)
	if err != nil {
		return nil, err
	}
	return o.objectConvert(data)
}

// Aggregate documents in the collection.
func (o *ArticleCategoryLogCollection) Aggregate(ctx context.Context, filter, projection, sort interface{}, offset, limit int64) ([]ArticleCategory, error) {
	data, err := o.dal.Aggregate(ctx, o.collection, filter, projection, sort, offset, limit)
	if err != nil {
		return nil, err
	}
	return o.dataConvert(data)
}

// Aggregate common documents in the collection.
func (o *ArticleCategoryLogCollection) AggregateCommon(ctx context.Context, pipeline interface{}) ([]ArticleCategory, error) {
	data, err := o.dal.AggregateCommon(ctx, o.collection, pipeline)
	if err != nil {
		return nil, err
	}
	return o.dataConvert(data)
}

// Aggregate common documents in the collection.
func (o *ArticleCategoryLogCollection) AggregateRaw(ctx context.Context, pipeline interface{}) ([]interface{}, error) {
	data, err := o.dal.AggregateCommon(ctx, o.collection, pipeline)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Count documents in the collection.
func (o *ArticleCategoryLogCollection) Counts(ctx context.Context, filter interface{}) (int64, error) {
	return o.dal.Counts(ctx, o.collection, filter)
}
