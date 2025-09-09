package main

import (
	"Project_5/main/controller"
	"fmt"
	"net/http"
	"os" // os is used to get environment variables
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

const testToken = "my-static-test-token"

// Get secret key from environment variable for security
var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// ACCESS TOKEN
func generateToken(username string) (string, error) {
	if len(jwtSecretKey) == 0 {
		return "", fmt.Errorf("JWT secret key not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 1 day expiration
	})

	return token.SignedString(jwtSecretKey)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 👉 Если токен совпадает с тестовым, пускаем
		if tokenStr == testToken {
			c.Set("username", "test-user")
			c.Next()
			return
		}

		// иначе обычная проверка JWT
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} //success

func main() {
	controller.ConnectDB()
	r := gin.Default()
	//	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	api := r.Group("/api")
	apiAuth := api.Group("/")
	apiAuth.POST("/handleVer", controller.HandleVerification)
	apiAuth.POST("/signup", controller.HandleSignUp)
	//apiAuth.POST("/depts", getDebts)
	apiAuth.Use(AuthMiddleware())
	{

		apiAuth.GET("/loans/getAllLoans", controller.GetLoansUsers)
		apiAuth.GET("/user", controller.GetUser)
		apiAuth.POST("/depts/addDept", controller.AddDept)
		apiAuth.GET("/reports", controller.GetMonthSummary)
		apiAuth.POST("/depts/remDept", controller.RemoveDept)
		apiAuth.GET("/reports/getHistory", controller.GetHistory)
	}

	r.Run("0.0.0.0:3030")
}
