package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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