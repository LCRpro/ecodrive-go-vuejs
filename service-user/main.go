package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
	"github.com/gin-gonic/gin"
)

var db *gorm.DB

func main() {
    dsn := "root:root@tcp(127.0.0.1:3306)/goecodrive?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    db.AutoMigrate(&User{})

	 r := gin.Default()

    r.POST("/users", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(400, gin.H{"error": "invalid body"})
            return
        }
        db.Create(&user)
        c.JSON(201, user)
    })

    r.GET("/users/:id", func(c *gin.Context) {
        var user User
        if err := db.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(404, gin.H{"error": "not found"})
            return
        }
        c.JSON(200, user)
    })

    r.Run(":8002")
}