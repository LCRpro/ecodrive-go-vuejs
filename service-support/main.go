package main

import (
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_DSN manquante")
	}

	if strings.HasPrefix(dsn, "mysql://") {
		parsedUrl, err := url.Parse(dsn)
		if err != nil {
			panic("Erreur parsing DATABASE_URL: " + err.Error())
		}
		user := parsedUrl.User.Username()
		pass, _ := parsedUrl.User.Password()
		host := parsedUrl.Host
		dbname := strings.TrimPrefix(parsedUrl.Path, "/")
		params := parsedUrl.RawQuery
		// Format final pour GORM
		dsn = user + ":" + pass + "@tcp(" + host + ")/" + dbname
		if params != "" {
			dsn += "?" + params
		}
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Erreur connexion BDD: " + err.Error())
	}
	db.AutoMigrate(&SupportTicket{})

	r := gin.Default()
frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
if frontendOrigin == "" {
	frontendOrigin = "http://localhost:5173" 
}

r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{frontendOrigin},
	AllowMethods:     []string{"GET", "POST", "PATCH"},
	AllowHeaders:     []string{"Content-Type"},
	AllowCredentials: true,
}))

	r.POST("/support", CreateSupportTicket)
	r.GET("/support/user/:google_id", ListUserTickets)
	r.GET("/support/all", ListAllTickets)
	r.PATCH("/support/reply/:id", ReplySupportTicket)

	r.Run(":8007")
}
