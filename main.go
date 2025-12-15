package main

import (
	"backend-path/configs"
	"backend-path/utils"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using env from machine")
	}

	app := fiberConfig()

	utils.ZapLogger(os.Getenv("APP_ENV"))

	configs.Setup()
	// repository.DB = configs.DB

	//middlewares.Setup(app)

	// routes.Setup(app)

	port := os.Getenv("APP_PORT")
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func fiberConfig() *fiber.App {
	maxBody, _ := strconv.Atoi(os.Getenv("APP_MAX_BODY"))
	if maxBody == 0 {
		maxBody = 4
	}

	app := fiber.New(fiber.Config{
		ReadBufferSize: maxBody * 1024,
		BodyLimit: maxBody * 1024,
	})

	return app 
}