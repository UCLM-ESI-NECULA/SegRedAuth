package main

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"seg-red-auth/internal/app/config"
)

func init() {
	_ = godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")
	certs := os.Getenv("CERTS_FOLDER")
	app := config.SetupRouter()
	err := app.RunTLS(":"+port, filepath.Join(certs, "auth.crt"), filepath.Join(certs, "auth.key"))
	if err != nil {
		panic(err)
	}
}
