package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/keyloom/web-api/controllers"
	docs "github.com/keyloom/web-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Gin Swagger setup
	e := gin.Default()
	docs.SwaggerInfo.BasePath = "/"
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Controller registration
	(&controllers.UserController{}).RegisterRoutes(e)
	(&controllers.TokenController{}).RegisterRoutes(e)
	(&controllers.AudienceController{}).RegisterRoutes(e)

	e.Run(":8080")
}
