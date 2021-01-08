package model

import (
	"context"
	"fmt"
	"golang-training/mongodb/config"
	"golang-training/mongodb/dal"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getDBSingleNodeConfig(t *testing.T) *config.DatabaseConfig {
	config := &config.DatabaseConfig{
		ConnectionURI: "mongodb://root:12345@localhost:27017/?authSource=go_mongo",
		Database:      "go_mongo",
		// Username:         "root",
		// Password:         "12345",
	}
	return config
}

func getDBConfig(t *testing.T) *config.DatabaseConfig {
	// mongodb://mongo1:9042,mongo2:9142,mongo3:9242/?replicaSet=rs0&readPreference=nearest
	// mongodb://root:12345@mongo1:9042,mongo2:9142,mongo3:9242/?replicaSet=rs0&readPreference=nearest&authSource=go_mongo
	config := &config.DatabaseConfig{
		// ConnectionURI: "mongodb://mongo1:9042,mongo2:9142,mongo3:9242/?replicaSet=rs0&readPreference=nearest",
		ConnectionURI: "mongodb://root:12345@mongo1:9042,mongo2:9142,mongo3:9242/?replicaSet=rs0&readPreference=nearest&authSource=go_mongo",
		Database:      "go_mongo",
		// Username:         "root",
		// Password:         "12345",
	}
	return config
}

func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}
	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}

func containsKey(doc bson.Raw, key ...string) bool {
	_, err := doc.LookupErr(key...)
	if err != nil {
		return false
	}
	return true
}

func parseDate(t *testing.T, dateString string) time.Time {
	rfc3339MilliLayout := "2006-01-02T15:04:05.999Z07:00" // layout defined with Go reference time
	parsedDate, err := time.Parse(rfc3339MilliLayout, dateString)

	require.NoError(t, err)
	return parsedDate
}

func TestNewArticleCategoryCollection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	_, err = NewArticleCategoryCollection(dal)
	require.NoError(t, err)
}

func TestArticleCategoryCollection_dataConvert(t *testing.T) {
	type args struct {
		listDocument []interface{}
	}
	tests := []struct {
		name     string
		o        *ArticleCategoryCollection
		args     args
		wantData []ArticleCategory
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := tt.o.dataConvert(tt.args.listDocument)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleCategoryCollection.dataConvert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ArticleCategoryCollection.dataConvert() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestArticleCategoryCollection_objectConvert(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		o       *ArticleCategoryCollection
		args    args
		want    *ArticleCategory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.objectConvert(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleCategoryCollection.objectConvert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleCategoryCollection.objectConvert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleCategoryCollection_Insert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	article_categories := []ArticleCategory{
		{
			SourceID: "lu",
			Title: []Title{
				{
					Text: "Title",
					Lang: "en",
				},
			},
			Description: []Description{
				{
					Text: "Desc",
					Lang: "en",
				},
			},
		},
		{
			SourceID: "nam",
			Title: []Title{
				{
					Text: "Tiêu đề",
					Lang: "vi",
				},
			},
			Description: []Description{
				{
					Text: "Mô tả",
					Lang: "vi",
				},
			},
		},
	}
	{
		// insert one
		result, err := coll.Insert(ctx, &article_categories[0])
		require.NoError(t, err)
		require.NotNil(t, result)
	}
	{
		// insert many
		result, err := coll.InsertMany(ctx, article_categories)
		require.NoError(t, err)
		require.Len(t, result, 2)
	}
}

func TestArticleCategoryCollection_Update(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	{
		// update one
		filter := bson.M{
			"source_id": "nam",
		}
		result, err := coll.FindOne(ctx, filter, nil)
		require.NoError(t, err)

		data := bson.M{
			"$set": bson.M{
				"running_status": RUNNING,
			},
			"$push": bson.M{
				"title": bson.M{
					"text": "newnew",
					"lang": "en",
				},
			},
		}
		filterByID := bson.D{primitive.E{Key: "_id", Value: result.ID}}
		err = coll.Update(ctx, filterByID, data, false)
		require.NoError(t, err)

		resultUp, err := coll.FindByID(ctx, result.ID, nil)
		require.NoError(t, err)
		require.Len(t, resultUp.Title, len(result.Title)+1)
		require.Equal(t, RUNNING, resultUp.RunningStatus)
	}
	{
		// update many
		filter := bson.M{
			"source_id": "lu",
		}
		data := bson.M{
			"$set": bson.M{
				"modified": "lux",
			},
		}
		err := coll.UpdateMany(ctx, filter, data, false)
		require.NoError(t, err)

		result, err := coll.Find(ctx, filter, nil, nil, 0, 10)
		require.NoError(t, err)
		for _, item := range result {
			require.Equal(t, "lux", item.Modified)
		}
	}
}

func TestArticleCategoryCollection_Find(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	var id primitive.ObjectID
	{
		// find all
		result, err := coll.Find(ctx, bson.D{}, nil, nil, 0, 10)
		require.NoError(t, err)
		require.Len(t, result, 3)
		id = result[0].ID
	}
	{
		// find with filter, projection, sort
		filter := bson.M{"source_id": "lu"}
		projection := bson.D{
			{Key: "_id", Value: 1},
			{Key: "status", Value: 1},
			{Key: "title", Value: 1},
		}
		sort := bson.D{
			{"status", 1},
		}
		result, err := coll.Find(ctx, filter, projection, sort, 0, 10)
		require.NoError(t, err)
		require.Len(t, result, 2)
		lastStatus := result[0].Status
		for _, item := range result {
			require.NotNil(t, item.ID)
			require.NotNil(t, item.Status)
			require.NotNil(t, item.Title)
			require.Nil(t, item.Description)
			require.LessOrEqual(t, lastStatus, item.Status)
			lastStatus = item.Status
		}
	}
	{
		// find by id
		result, err := coll.FindByID(ctx, id.Hex(), nil)
		require.NoError(t, err)
		require.Equal(t, id, result.ID)
	}
	{
		// find one
		filter := bson.M{"source_id": "lu"}
		result, err := coll.FindOne(ctx, filter, nil)
		require.NoError(t, err)
		require.Equal(t, "lu", result.SourceID)
	}
	{
		// find distinct
		result, err := coll.Distinct(ctx, "source_id", bson.D{})
		require.NoError(t, err)
		require.Len(t, result, 2)
	}
	{
		// filter with time
		initDate := time.Now().UTC()
		// initDate, err := time.Parse("02/01/2006 15:04:05", now)
		// require.NoError(t, err)
		filter := bson.D{
			{Key: "updatedAt", Value: bson.D{
				{Key: "$gt", Value: initDate},
			}},
		}
		result, err := coll.Find(ctx, filter, nil, nil, 0, 10)
		require.NoError(t, err)
		require.Len(t, result, 0)

		dateBefore := time.Now().Add(-24 * time.Hour).UTC()
		filter = bson.D{
			{"updatedAt", bson.D{
				{"$gt", dateBefore},
			}},
		}
		result, err = coll.Find(ctx, filter, nil, nil, 0, 10)
		require.NoError(t, err)
		require.Len(t, result, 3)
	}
}

func TestArticleCategoryCollection_Aggregate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	{
		result, err := coll.Aggregate(ctx, bson.D{}, nil, nil, 0, 10)
		require.NoError(t, err)
		require.Len(t, result, 3)
	}
	{
		// return struct
		pipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"source_id", "lu"},
				}},
			},
			{
				{"$project", bson.D{
					{"dayOfWeek", "$_id.day"},
					{"_id", 1},
				}},
			},
			{
				{"$sort", bson.D{
					{"updatedAt", -1},
				}},
			},
		}
		result, err := coll.AggregateCommon(ctx, pipeline)
		require.NoError(t, err)
		require.Len(t, result, 2)
	}
	{
		// return []interface{}
		pipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"source_id", "lu"},
				}},
			},
			{
				{"$project", bson.D{
					{"dayOfWeek", bson.D{
						{"$dayOfWeek", "$updatedAt"},
					}},
					{"_id", 1},
					{"updatedAt", 1},
				}},
			},
			{
				{"$sort", bson.D{
					{"updatedAt", -1},
				}},
			},
		}
		result, err := coll.AggregateRaw(ctx, pipeline)
		require.NoError(t, err)
		require.Len(t, result, 2)
		for _, item := range result {
			require.Contains(t, item, "_id")
			require.Contains(t, item, "dayOfWeek")
		}
	}
	{
		// with json pipeline
		pipeline := `[
			{
				"$group": {
					"_id": "$source_id",
					"objectIDs": {
						"$addToSet": "$_id"
					},
					"count": {
						"$sum": 1
					}
				}
			},
			{
				"$sort": {
					"count": 1
				}
			}
		]`
		result, err := coll.AggregateRaw(ctx, pipeline)
		require.NoError(t, err)
		// log with -v flag
		t.Logf("Aggregate with json pipeline: %#v", result)
		require.Len(t, result, 2)
		var lastCount int
		for i, item := range result {
			require.Contains(t, item, "_id")
			require.Contains(t, item, "count")
			require.Contains(t, item, "objectIDs")
			// count == len(objectIDs)
			item, ok := item.(map[string]interface{})
			require.True(t, ok)
			len, err := strconv.Atoi(fmt.Sprintf("%v", item["count"]))
			require.NoError(t, err)
			require.Len(t, item["objectIDs"], len)
			// sort count
			if i == 0 {
				lastCount = len
			}
			require.LessOrEqual(t, lastCount, len)
			lastCount = len
		}
	}
}

func TestArticleCategoryCollection_Counts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	{
		result, err := coll.Counts(ctx, bson.D{})
		require.NoError(t, err)
		require.Equal(t, int64(3), result)
	}
	{
		// count with filter
		// filter := bson.D{{Key: "source_id", Value: "lu"}}
		filter := bson.M{"source_id": "lu"}
		result, err := coll.Counts(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, int64(2), result)
	}
	{
		// count not found
		filter := bson.M{"source_id": "luxxx"}
		result, err := coll.Counts(ctx, filter)
		require.NoError(t, err)
		require.Equal(t, int64(0), result)
	}
}

// Transaction numbers are only allowed on a replica set member or mongos
func TestArticleCategoryCollection_Transaction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := getDBConfig(t)
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	require.NoError(t, err)
	defer dal.Disconnect(ctx)

	client, err := dal.GetClient(ctx)
	require.NoError(t, err)
	clientMongo, ok := client.(*mongo.Client)
	require.True(t, ok)

	coll, err := NewArticleCategoryCollection(dal)
	require.NoError(t, err)

	{
		// success
		filter := bson.M{
			"source_id": "nam",
		}
		result, err := coll.FindOne(ctx, filter, nil)
		require.NoError(t, err)

		session, err := clientMongo.StartSession()
		require.NoError(t, err)

		err = session.StartTransaction()
		require.NoError(t, err)

		err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
			// commands
			data := bson.M{
				"$set": bson.M{
					"running_status": RUNNING,
				},
				"$push": bson.M{
					"title": bson.M{
						"text": "transaction",
						"lang": "en",
					},
				},
			}
			filterByID := bson.D{primitive.E{Key: "_id", Value: result.ID}}
			err = coll.Update(sc, filterByID, data, false)
			require.NoError(t, err)

			if err = session.CommitTransaction(sc); err != nil {
				t.Fatal(err)
			}
			return nil
		})
		require.NoError(t, err)
		session.EndSession(ctx)

		resultUp, err := coll.FindByID(ctx, result.ID, nil)
		require.NoError(t, err)
		require.Len(t, resultUp.Title, len(result.Title)+1)
		require.Equal(t, RUNNING, resultUp.RunningStatus)
	}
	{
		// abort
		filter := bson.M{
			"source_id": "lu",
		}

		session, err := clientMongo.StartSession()
		require.NoError(t, err)

		err = session.StartTransaction()
		require.NoError(t, err)

		err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
			// commands
			data := bson.M{
				"$set": bson.M{
					"modified": "luxxx",
				},
			}
			err := coll.UpdateMany(sc, filter, data, false)
			require.NoError(t, err)

			// abort
			session.AbortTransaction(sc)
			return nil
		})
		require.NoError(t, err)
		session.EndSession(ctx)

		result, err := coll.Find(ctx, filter, nil, nil, 0, 10)
		require.NoError(t, err)
		for _, item := range result {
			require.Equal(t, "lux", item.Modified)
		}
	}
}
