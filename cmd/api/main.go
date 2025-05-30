package main

import (
	"fmt"
	"os"

	"github.com/gaspartv/encurtador-de-url/internal/configs"
	"github.com/gaspartv/encurtador-de-url/internal/routes"
	"github.com/joho/godotenv"
)

var (
	logger configs.Logger
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
		os.Exit(1)
	}

	logger = *configs.GetLogger("main")

	routes.InitRoutes()
}
