package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() {
  client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
    panic(err)
  }

  usersCollection := client.Database("testing-go").Collection("users")

  // insert a single document into a collection
  // create a bson.D object
  user := bson.D{{"fullName", "User 1"}, {"age", 30}}
  // insert the bson object using InsertOne()
  // result, err := usersCollection.InsertOne(context.TODO(), user)
  // check for errors in the insertion
  if err != nil {
    panic(err)
  }
  // display the id of the newly inserted object
  // fmt.Println(result.InsertedID)
  fmt.Println(user, usersCollection)
}