package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	UserId   primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

// --add user--
func addUser(c echo.Context, client *mongo.Client) error {

	//--get user details--
	username := c.QueryParam("username")
	email := c.QueryParam("email")
	password := c.QueryParam("password")

	coll := client.Database("UserDB").Collection("Users")
	newUser := User{Username: username, Email: email, Password: password}

	result, err := coll.InsertOne(context.TODO(), newUser)
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, fmt.Sprintf("Document inserted with ID: %v\n", result.InsertedID))
}

// --delete user--
func deleteUser(c echo.Context, client *mongo.Client) error {

	//--get username to delete--
	username := c.QueryParam("username")

	coll := client.Database("UserDB").Collection("Users")
	filter := bson.D{{"username", username}}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, fmt.Sprintf(" %v Document deleted with Username: %s\n", result.DeletedCount, username))
}

// --update username--
func updateUser(c echo.Context, client *mongo.Client) error {

	username := c.QueryParam("username")
	newname := c.QueryParam("newname")

	coll := client.Database("UserDB").Collection("Users")

	//--get userID--
	filteru := bson.D{{"username", username}}
	var resultu User
	var err = coll.FindOne(context.TODO(), filteru).Decode(&resultu)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.String(http.StatusOK, fmt.Sprint("no documents found"))
		}
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	id, _ := primitive.ObjectIDFromHex(resultu.UserId.String())
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"username", newname}}}}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, string(output))
}

// --get user--
func getUser(c echo.Context, client *mongo.Client) error {

	username := c.QueryParam("username")

	coll := client.Database("UserDB").Collection("Users")
	filter := bson.D{{"username", username}}

	var result User
	var err = coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.String(http.StatusOK, fmt.Sprint("no documents found"))
		}
		panic(err)
	}
	output, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, string(output))
}

func main() {

	//--connect to the mongodb cluster--
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := "mongodb+srv://chamiduhp48:52vlvzl8q9KPFNLx@cluster0.l3mrfxr.mongodb.net/?retryWrites=true&w=majority"
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//--echo--
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "MongoDB test with Go")
	})
	//--add a user--
	e.POST("/add", func(c echo.Context) error {
		return addUser(c, client)
	})

	//--change user details--
	e.PUT("/update", func(c echo.Context) error {
		return updateUser(c, client)
	})

	//--get a user--
	e.GET("/get", func(c echo.Context) error {
		return getUser(c, client)
	})

	//--delete a user--
	e.DELETE("/delete", func(c echo.Context) error {
		return deleteUser(c, client)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
