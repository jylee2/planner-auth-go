package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"api-go/database"
	"api-go/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
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

const CookieName = "jwt"
const MongoDB = "auth-go"
const Model = "users"

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14) // cost: 14

	user := bson.D{
		primitive.E{Key: "uuid", Value: uuid.New().String()},
		primitive.E{Key: "name", Value: data["name"]},
		primitive.E{Key: "username", Value: data["username"]},
		primitive.E{Key: "email", Value: data["email"]},
		primitive.E{Key: "password", Value: password},
	}

	client, _, _, err := database.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	usersCollection := client.Database(MongoDB).Collection(Model)

	filter := bson.D{
		primitive.E{Key: "username", Value: data["username"]},
	}

	// Find existing user
	cursor, _ := usersCollection.Find(context.TODO(), filter)
	// convert the cursor result to bson
	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		if result["username"] == data["username"] {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Username already exists.",
			})
		}
	}

	// insert a single document into a collection
	// create a bson.D object
	// user := bson.D{{"name", "User 1"}, {"email", "test@test.com"}}
	// insert the bson object using InsertOne()
	_, insertErr := usersCollection.InsertOne(context.TODO(), user)
	// check for errors in the insertion
	if insertErr != nil {
		panic(insertErr)
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	client, _, _, err := database.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	usersCollection := client.Database(MongoDB).Collection(Model)

	filter := bson.D{
		primitive.E{Key: "username", Value: data["username"]},
	}

	// retrieving the first document that match the filter
	// var user bson.M
	var user models.User
	// check for errors in the finding
	if err = usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		fmt.Println("--------Login usersCollection.FindOne Error: ", err)
		// panic(err)
		// User not found.
		// c.Status(fiber.StatusNotFound)
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Incorrect username or password.",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Incorrect username or password.",
		})
	}

	expiryInHours, err := time.ParseDuration(os.Getenv("JWT_EXPIRY_IN_HOURS"))
	if err != nil {
		panic(err)
	}

	jwtExpiry := time.Now().Add(time.Hour * expiryInHours)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Uuid,
		ExpiresAt: jwtExpiry.Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("AUTH_SECRET")))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"success": false,
			"message": "Could not login.",
		})
	}

	cookie := fiber.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  jwtExpiry,
		HTTPOnly: true,
	}

	// Set the cookie
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"success": true,
		"jwt":     token,
		// "userUuid": user.Uuid,
	})
}

func GetUserFromCookie(c *fiber.Ctx) error {
	cookie := c.Cookies(CookieName)

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_SECRET")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated.",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims) // convert to StandardClaims

	client, _, _, err := database.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}
	usersCollection := client.Database(MongoDB).Collection(Model)
	filter := bson.D{
		primitive.E{Key: "uuid", Value: claims.Issuer},
	}

	var user models.User
	// Find user by uuid
	if err = usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		// panic(err)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found.",
		})
	}

	return c.JSON(user)
}

func GetUserByUuid(c *fiber.Ctx) error {
	// cookie := c.Params(CookieName)

	// _, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(os.Getenv("AUTH_SECRET")), nil
	// })

	// fmt.Println("--------err: ", err)
	// if err != nil {
	// 	c.Status(fiber.StatusUnauthorized)
	// 	return c.JSON(fiber.Map{
	// 		"message": "Unauthenticated.",
	// 	})
	// }

	client, _, _, err := database.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}
	usersCollection := client.Database(MongoDB).Collection(Model)
	filter := bson.D{
		primitive.E{Key: "uuid", Value: c.Params("uuid")},
	}

	var user models.User
	// Find user by uuid
	if err = usersCollection.FindOne(context.TODO(), filter).Decode(&user); err != nil {
		// panic(err)
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found.",
		})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // expires 1 hour ago
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"success": true,
	})
}
