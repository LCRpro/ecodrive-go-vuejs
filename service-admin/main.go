package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var userServiceURL string

func main() {
	godotenv.Load()

	userServiceURL = os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8002"
	}

	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET manquant dans .env")
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(AuthRequired())

	admin := r.Group("/admin")
	admin.Use(AdminRequired())

	{
		admin.GET("/users", ListUsers)
		admin.PATCH("/users/:id/accept-driver", AcceptDriver)
		admin.GET("/driver-requests", ListDriverRequests)
		admin.PATCH("/driver-requests/:id", HandleDriverRequest)
	}
	r.POST("/admin/driver-requests", CreateDriverRequest)

	r.Run(":8003")
}
