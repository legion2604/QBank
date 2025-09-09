package main

import (
	"Project_5/main/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	controller.ConnectDB()
	r := gin.Default()
	//	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(controller.CORSMiddleware())

	api := r.Group("/api")
	apiAuth := api.Group("/")
	apiAuth.POST("/handleVer", controller.HandleVerification)
	apiAuth.POST("/signup", controller.HandleSignUp)
	//apiAuth.POST("/depts", getDebts)
	apiAuth.Use(controller.AuthMiddleware())
	{

		apiAuth.GET("/loans/getAllLoans", controller.GetLoansUsers)
		apiAuth.GET("/user", controller.GetUser)
		apiAuth.POST("/depts/addDept", controller.AddDept)
		apiAuth.GET("/reports", controller.GetMonthSummary)
		apiAuth.POST("/depts/remDept", controller.RemoveDept)
		apiAuth.GET("/reports/getHistory", controller.GetHistory)
	}

	r.Run("localhost:8080")
}
