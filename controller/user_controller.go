package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chamidu48/go_with_mongoDB/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	session *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	return &UserController{s}
}

// --get user details--
func (uc UserController) GetUser(c echo.Context) error {
	username := c.QueryParam("username")

	coll := uc.session.Database("UserDB").Collection("Users")
	filter := bson.D{{"username", username}}

	var result models.UserR
	var err = coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			//--query did not match any documents--
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

// --add a user to db--
func (uc UserController) AddUser(c echo.Context) error {
	//--get user details from JSON payload--
	var newUser models.User
	if err := c.Bind(&newUser); err != nil {
		return err
	}

	coll := uc.session.Database("UserDB").Collection("Users")

	result, err := coll.InsertOne(context.TODO(), newUser)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "User added successfully!",
		"userId":   result.InsertedID,
		"username": newUser.Username,
		"email":    newUser.Email,
		"password": newUser.Password,
	})
}

// --update username--
func (uc UserController) UpdateUserName(c echo.Context) error {
	username := c.QueryParam("username")
	newname := c.QueryParam("newname")

	coll := uc.session.Database("UserDB").Collection("Users")

	//--get userID--
	filteru := bson.D{{"username", username}}
	var resultu models.UserR
	var err = coll.FindOne(context.TODO(), filteru).Decode(&resultu)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.String(http.StatusOK, fmt.Sprint("no documents found"))
		}
		panic(err)
	}

	id, err := primitive.ObjectIDFromHex(resultu.UserId.Hex())
	if err != nil {
		panic(err)
	}
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

// --delete a user--
func (uc UserController) DeleteUser(c echo.Context) error {
	//--get username to delete--
	username := c.QueryParam("username")

	coll := uc.session.Database("UserDB").Collection("Users")
	filter := bson.D{{"username", username}}

	result, err := coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, fmt.Sprintf(" %v Document deleted with Username: %s\n", result.DeletedCount, username))
}

// --check username--
func (uc UserController) CheckUser(c echo.Context) error {
	var userAuth models.UserAuth
	var user models.UserAuth

	if err := c.Bind(&userAuth); err != nil {
		return err
	}
	coll := uc.session.Database("UserDB").Collection("Users")
	filter := bson.D{{"username", userAuth.Username}}

	err := coll.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Not registered",
			})
		}
		if userAuth.Password != user.Password {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Incorrect password" + user.Password,
			})
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success" + user.Password,
	})
}

// --get all users--
func (uc UserController) GetAll(c echo.Context) error {
	// Access the "UserDB" database and "Users" collection
	coll := uc.session.Database("UserDB").Collection("Users")

	// Execute a find operation with an empty filter to retrieve all documents
	cur, err := coll.Find(context.Background(), bson.D{})
	if err != nil {
		panic(err) // If there's an error, panic (halt the execution)
	}

	// Iterate over the cursor and append the results to a slice
	var results []bson.M
	if err := cur.All(context.Background(), &results); err != nil {
		panic(err)
	}

	// Return the results as a response (assuming you want to return JSON)
	return c.JSON(http.StatusOK, results)
}
