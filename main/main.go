package main

// @title Project_5 API
// @version 1.0
// @description API для управления долгами
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"Project_5/main/controller"
	_ "Project_5/main/docs" // Обновленный импорт, так как docs находится в корне
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	controller.ConnectDB()
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(controller.CORSMiddleware())

	api := r.Group("/api")
	apiAuth := api.Group("/")
	apiAuth.POST("/handleVer", controller.HandleVerification)
	apiAuth.POST("/signup", controller.HandleSignUp)
	apiAuth.Use(controller.AuthMiddleware())
	{
		apiAuth.GET("/loans/getAllLoans", controller.GetLoansUsers)
		apiAuth.GET("/user", controller.GetUser)
		apiAuth.POST("/deps/addDept", controller.AddDept)
		apiAuth.GET("/reports", controller.GetMonthSummary)
		apiAuth.POST("/deps/remDept", controller.RemoveDept)
		apiAuth.GET("/reports/getHistory", controller.GetHistory)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("localhost:8080")
}
