package controllers

import (
	"context"
	"fmt"
	"time"

	"api-go/database"
	"api-go/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type E struct {
	Key   string
	Value interface{}
}

const SecretKey = "secret"
const CookieName = "jwt"
const MongoUri = "mongodb://localhost:27017"
const MongoDB = "testing-go"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) // cost: 14

	// user := models.User{
	// 	Uuid: uuid.New(),
	// 	Name: data["name"],
	// 	Email: data["email"],
	// 	Password: password,
	// }
	// filter := bson.D{primitive.E{Key: "autorefid", Value: "100"}}
  // user := bson.D{{"name", data["name"]}, {"email", data["email"]}, {"password", password}}

	// fmt.Println("--------reflect.TypeOf(uuid.New()): ", reflect.TypeOf(uuid.New()))
	// fmt.Println("--------reflect.TypeOf(uuid.New().String()): ", reflect.TypeOf(uuid.New().String()))
  user := bson.D{
		primitive.E{Key: "uuid", Value: uuid.New().String()},
		primitive.E{Key: "name", Value: data["name"]},
		primitive.E{Key: "email", Value: data["email"]},
		primitive.E{Key: "password", Value: password},
	}
  fmt.Println("--------user: ", user)

	client, _, _, err := database.Connect(MongoUri)
  if err != nil {
		panic(err)
  }

  usersCollection := client.Database(MongoDB).Collection("users")

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

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	client, _, _, err := database.Connect(MongoUri)
  if err != nil {
		panic(err)
  }

  usersCollection := client.Database(MongoDB).Collection("users")

	// create a search filer
	// filter := bson.D{
	// 	{"$and",
	// 		bson.A{
	// 			bson.D{
	// 				{"age", bson.D{{"$gt", 25}}},
	// 			},
	// 		},
	// 	},
	// }
	// filter := bson.D{
	// 	primitive.E{Key: "uuid", Value: uuid.New()},
	// 	primitive.E{Key: "name", Value: data["name"]},
	// 	primitive.E{Key: "email", Value: data["email"]},
	// 	primitive.E{Key: "password", Value: password},
	// }
	filter := bson.D{
		primitive.E{Key: "email", Value: data["email"]},
	}
  fmt.Println("--------filter: ", filter)

	// retrieving the first document that match the filter
	// var user bson.M
	var user models.User
	// check for errors in the finding
	if err = usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		fmt.Println("--------usersCollection.FindOne err: ", err)
		// panic(err)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found.",
		})
	}

	// if user == nil {
	// 	c.Status(fiber.StatusNotFound)
	// 	return c.JSON(fiber.Map{
	// 		"message": "User not found.",
	// 	})
  // }

	// display the document retrieved
	fmt.Println("--------user: ", user)

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password.",
		})
	}

	jwtExpiry := time.Now().Add(time.Hour * 24) // 24 hours
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		// Issuer: strconv.FormatUint(uint64(user.Uuid), 10),
		Issuer: user.Uuid,
		ExpiresAt: jwtExpiry.Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login.",
		})
	}

	cookie := fiber.Cookie{
		Name: CookieName,
		Value: token,
		Expires: jwtExpiry,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func GetUserFromCookie(c *fiber.Ctx) error {
	cookie := c.Cookies(CookieName)

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "The user is not logged in.",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) // convert to StandardClaims
	fmt.Println("--------claims: ", claims)

	client, _, _, err := database.Connect(MongoUri)
  if err != nil {
		panic(err)
  }

  usersCollection := client.Database(MongoDB).Collection("users")
	filter := bson.D{
		primitive.E{Key: "uuid", Value: claims.Issuer},
	}

	var user models.User
	// check for errors in the finding
	if err = usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		// panic(err)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found.",
		})
	}
	fmt.Println("--------user: ", user)

	return c.JSON(user)
}