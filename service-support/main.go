package main

import (
    "os"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var db *gorm.DB

func main() {
    godotenv.Load()

    dsn := os.Getenv("DATABASE_DSN") 
    if dsn == "" {
        panic("DATABASE_DSN manquante")
    }

    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Erreur connexion BDD: " + err.Error())
    }

    r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
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