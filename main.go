package main

import (
	"backend-path/app/middlewares"
	"backend-path/app/repository"
	"backend-path/configs"
	"backend-path/database/seeders"
	"backend-path/routes"
	"backend-path/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	repository.DB = configs.DB

	
	middlewares.Setup(app)
	
	routes.Setup(app)
	argsListener()

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
		BodyLimit:      maxBody * 1024,
	})

	return app
}

func argsListener() {
	homeDir, _ := os.UserHomeDir()
	sqlMigrate := homeDir + "/go/bin/sql-migrate"

	for _, arg := range os.Args {
		if arg == "--rollback" {
			utils.Logger.Info("✅ down all migration")
			cmd := exec.Command(sqlMigrate, "down", "-limit=0")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			utils.Logger.Info("✅ up all migration")
			cmd = exec.Command(sqlMigrate, "up", "-limit=0")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		if arg == "--seed" {
			runner := seeders.All(configs.DB)
			if err := runner.Run(); err != nil {
				panic(err)
			}
		}
	}
}
