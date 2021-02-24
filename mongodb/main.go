package main

import (
	"context"
	"fmt"
	"golang-training/mongodb/config"
	"golang-training/mongodb/dal"
	"golang-training/mongodb/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	config := &config.DatabaseConfig{
		ConnectionURI: "mongodb://root:12345@localhost:27017/?authMechanism=SCRAM-SHA-256&authSource=go_mongo",
		Database:      "go_mongo",
		// Username:         "root",
		// Password:         "12345",
	}
	dal, err := dal.NewDataAccessLayerMongoDB(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer dal.Disconnect(ctx)

	// article_category := model.ArticleCategory{
	// 	ID:       primitive.NewObjectID(),
	// 	SourceID: "lu",
	// 	Title: []model.Title{
	// 		{
	// 			Text: "New",
	// 			Lang: "en",
	// 		},
	// 	},
	// 	Description: []model.Description{
	// 		{
	// 			Text: "New",
	// 			Lang: "en",
	// 		},
	// 	},
	// }

	// using direct dal
	// collectionName := "article_category"
	// result, err := dal.Insert(ctx, collectionName, article_category)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("result", result)

	// article_category collection
	aModel, err := model.NewArticleCategoryModel(dal)
	if err != nil {
		log.Fatal(err)
	}
	// result, err = aModel.Insert(ctx, &article_category)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("result2", result)

	projection := bson.D{{}}
	data, err := aModel.FindByID(ctx, "5ff691d1b6c552a7bd1d94ec", projection)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("findById", data)
	fmt.Println("")

	data1, err := aModel.Find(ctx, bson.D{}, nil, nil, 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("find", data1)

	fmt.Println("Done")
}
