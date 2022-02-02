package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func Connect() {
//   client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
//   if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
//     panic(err)
//   }

//   // usersCollection := client.Database("testing-go").Collection("users")

//   // insert a single document into a collection
//   // create a bson.D object
//   // user := bson.D{{"fullName", "User 1"}, {"age", 30}}
//   // insert the bson object using InsertOne()
//   // result, err := usersCollection.InsertOne(context.TODO(), user)
//   // check for errors in the insertion
//   if err != nil {
//     panic(err)
//   }
//   // display the id of the newly inserted object
//   // fmt.Println(result.InsertedID)
//   // fmt.Println(user, usersCollection)
//   fmt.Println("--------Connected to MongoDB.")
// }

// ==========

// This is a user defined method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and
// resource associated with it.
 
func Connect(uri string)(*mongo.Client, context.Context,
  context.CancelFunc, error) {
   
  // ctx will be used to set deadline for process, here
  // deadline will of 30 seconds.
  ctx, cancel := context.WithTimeout(context.Background(),
    30 * time.Second)

  // mongo.Connect return mongo.Client method
  client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
  return client, ctx, cancel, err
}