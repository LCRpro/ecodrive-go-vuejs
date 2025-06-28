package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var db *gorm.DB

func main() {
    dsn := "root:root@tcp(127.0.0.1:8889)/ecodrive?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Erreur connexion BDD: " + err.Error())
    }
    db.AutoMigrate(&User{})

    r := gin.Default()
    r.POST("/users", CreateUser)
    r.GET("/users/:id", GetUser)

    r.Run(":8002")
}