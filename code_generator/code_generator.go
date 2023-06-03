package code_generator

import (
	"net/http"
	"os"
	"text/template"

	"github.com/chamidu48/go_with_mongoDB/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type CodeGenerator struct {
	session *mongo.Client
}

func NewCodeGenerator(s *mongo.Client) *CodeGenerator {
	return &CodeGenerator{s}
}

// --code generarion--
func (uc CodeGenerator) Generate(c echo.Context) error {
	var templateModel models.TemplateModel
	if err := c.Bind(&templateModel); err != nil {
		return err
	}

	mainTemplatePath := "templates/main.txt"
	modelTemplatePath := "templates/model.txt"
	controllerTemplatePath := "templates/controller.txt"
	envTempPath := "templates/env.txt"

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

	envTemp, err := template.ParseFiles(envTempPath)
	if err != nil {
		return err
	}

	// Execute the templates and write to files
	mainOutputPath := "code_generated/main_generated.go"
	modelOutputPath := "code_generated/generated_models/model_generated.go"
	controllerOutputPath := "code_generated/generated_controller/controller_generated.go"
	envOutputPath := "code_generated/.env"

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

	envFile, err := os.Create(envOutputPath)
	if err != nil {
		return err
	}
	defer envFile.Close()

	err = envTemp.Execute(envFile, templateModel)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Code generation completed")
}
