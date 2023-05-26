package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"os"
	"text/template"

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

// --code generarion--
func (uc UserController) Generate(c echo.Context) error {
	var template_model models.TemplateModel
	if err := c.Bind(&template_model); err != nil {
		return err
	}

	main_template_path := "D:/Go/Go with mongoDB/templates/main.txt"
	model_template_path := "D:/Go/Go with mongoDB/templates/model.txt"
	controller_template_path := "D:/Go/Go with mongoDB/templates/controller.txt"

	//--init the templates--
	main_temp, err := template.New("main.txt").ParseFiles(main_template_path)
	if err != nil {
		return err
	}

	model_temp, err := template.New("model.txt").ParseFiles(model_template_path)
	if err != nil {
		return err
	}

	controller_temp, err := template.New("controller.txt").ParseFiles(controller_template_path)
	if err != nil {
		return err
	}

	//--executing the templates--
	err1 := main_temp.Execute(os.Stdout, template_model)
	if err1 != nil {
		return err1
	}

	fmt.Println("\n")

	err2 := model_temp.Execute(os.Stdout, template_model)
	if err2 != nil {
		return err1
	}

	fmt.Println("\n")

	err3 := controller_temp.Execute(os.Stdout, template_model)
	if err3 != nil {
		return err1
	}

	return c.String(http.StatusOK, string("Done"))
}
