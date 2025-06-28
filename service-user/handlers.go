package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
)

func FindOrCreateUser(c *gin.Context) {
    var input struct {
        GoogleID string `json:"google_id"`
        Email    string `json:"email"`
        Name     string `json:"name"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    var user User
    if err := db.Where("google_id = ?", input.GoogleID).First(&user).Error; err == gorm.ErrRecordNotFound {
        user = User{
            GoogleID: input.GoogleID,
            Email:    input.Email,
            Name:     input.Name,
            Roles:    `["ROLE_PASSAGER"]`, 
        }
        db.Create(&user)
    }
    c.JSON(200, user)
}

func GetUserByID(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    c.JSON(200, user)
}