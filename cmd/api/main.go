package main

import (
	"github.com/gaspartv/encurtador-de-url/internal/routes"
)

type formA struct {
	Foo string `json:"foo" xml:"foo" binding:"required"`
}

func main() {
	routes.InitRoutes()
}
