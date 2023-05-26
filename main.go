package main

import (
	"context"
	"log"

	"net/http"

	"github.com/chamidu48/go_with_mongoDB/controller"
	"github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	//--init the controller--
	uc := controller.NewUserController(client)

	//--echo--
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "MongoDB test with Go")
	})
	//--add a user--
	e.POST("/add", func(c echo.Context) error {
		return uc.AddUser(c)
	})

	//--change user details--
	e.PUT("/update", func(c echo.Context) error {
		return uc.UpdateUserName(c)
	})

	//--get a user--
	e.GET("/get", func(c echo.Context) error {
		return uc.GetUser(c)
	})

	//--delete a user--
	e.DELETE("/delete", func(c echo.Context) error {
		return uc.DeleteUser(c)
	})

	//--code generation using text/templates--
	e.GET("/generate", func(c echo.Context) error {
		return uc.Generate(c)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
