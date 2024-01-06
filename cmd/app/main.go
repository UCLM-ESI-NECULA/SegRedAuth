package main

import (
	"github.com/joho/godotenv"
	"os"
	"seg-red-auth/internal/app/config"
)

func init() {
	_ = godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")
	app := config.SetupRouter()
	err := app.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
