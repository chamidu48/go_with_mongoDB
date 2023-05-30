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
	var templateModel models.TemplateModel
	if err := c.Bind(&templateModel); err != nil {
		return err
	}

	mainTemplatePath := "D:/Go/Go with mongoDB/templates/main.txt"
	modelTemplatePath := "D:/Go/Go with mongoDB/templates/model.txt"
	controllerTemplatePath := "D:/Go/Go with mongoDB/templates/controller.txt"

	// Initialize the templates
	mainTemp, err := template.ParseFiles(mainTemplatePath)
	if err != nil {
		return err
	}

	modelTemp, err := template.ParseFiles(modelTemplatePath)
	if err != nil {
		return err
	}

	controllerTemp, err := template.ParseFiles(controllerTemplatePath)
	if err != nil {
		return err
	}

	// Execute the templates and write to files
	mainOutputPath := "D:/Go/Go with mongoDB/generated_temp/main_generated.txt"
	modelOutputPath := "D:/Go/Go with mongoDB/generated_temp/model_generated.txt"
	controllerOutputPath := "D:/Go/Go with mongoDB/generated_temp/controller_generated.txt"

	mainFile, err := os.Create(mainOutputPath)
	if err != nil {
		return err
	}
	defer mainFile.Close()

	err = mainTemp.Execute(mainFile, templateModel)
	if err != nil {
		return err
	}

	modelFile, err := os.Create(modelOutputPath)
	if err != nil {
		return err
	}
	defer modelFile.Close()

	err = modelTemp.Execute(modelFile, templateModel)
	if err != nil {
		return err
	}

	controllerFile, err := os.Create(controllerOutputPath)
	if err != nil {
		return err
	}
	defer controllerFile.Close()

	err = controllerTemp.Execute(controllerFile, templateModel)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Code generation completed")
}
