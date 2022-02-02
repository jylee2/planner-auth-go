package controllers

import (
	"context"
	"fmt"

	"api-go/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) // cost: 14

	// user := models.User{
	// 	Name: data["name"],
	// 	Email: data["email"],
	// 	Password: password,
	// }
  user := bson.D{{"name", data["name"]}, {"email", data["email"]}, {"password", password}}
  fmt.Println("--------user: ", user)

	client, _, _, err := database.Connect("mongodb://localhost:27017")
  if err != nil {
		panic(err)
  }

  usersCollection := client.Database("testing-go").Collection("users")

  // insert a single document into a collection
  // create a bson.D object
  // user := bson.D{{"fullName", "User 1"}, {"age", 30}}
  // insert the bson object using InsertOne()
  insertRes, insertErr := usersCollection.InsertOne(context.TODO(), user)
  // check for errors in the insertion
  if insertErr != nil {
    panic(insertErr)
  }
  // display the id of the newly inserted object
  // fmt.Println(insertRes.InsertedID)
  // fmt.Println(user, usersCollection)
  fmt.Println("--------insertRes: ", insertRes)

	return c.JSON(user)
}