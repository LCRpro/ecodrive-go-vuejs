package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	if os.Getenv("JWT_SECRET") == "" {
		panic("JWT_SECRET manquant dans .env")
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ecodrive.liamcariou.fr"},
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
