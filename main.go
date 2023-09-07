package main

import (
	"coba_01/docs"
	"coba_01/pkg/env"
	"coba_01/src/app/controller"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func init() {
	env.LoadEnvVariables()
}

// @title Implementasi-gin-gonic-crud-jwt API
// @version 0.0.1
// @description This is the API documentation for implementasi-gin-gonic-crud-jwt.
// @host localhost:5000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()
	baseApi := "/api/v1"
	docs.SwaggerInfo.BasePath = baseApi
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api1 := r.Group(baseApi)
	controller.NewUserController(api1)
	log.Fatal(r.Run(":5000"))
}
