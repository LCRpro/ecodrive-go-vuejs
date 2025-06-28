package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
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
        roles := []string{"ROLE_PASSAGER"}
        if input.Email == "liam2106cariou@gmail.com" {
            roles = append(roles, "ROLE_ADMIN")
        }
        rolesJSON, _ := json.Marshal(roles)
        user = User{
            GoogleID: input.GoogleID,
            Email:    input.Email,
            Name:     input.Name,
            Roles:    string(rolesJSON),
        }
        db.Create(&user)
    } else {
        var roles []string
        json.Unmarshal([]byte(user.Roles), &roles)
        updated := false
        if input.Email == "liam2106cariou@gmail.com" && !contains(roles, "ROLE_ADMIN") {
            roles = append(roles, "ROLE_ADMIN")
            updated = true
        }
        if updated {
            rolesJSON, _ := json.Marshal(roles)
            user.Roles = string(rolesJSON)
            db.Save(&user)
        }
    }
    c.JSON(200, user)
}

func contains(list []string, item string) bool {
    for _, v := range list {
        if v == item {
            return true
        }
    }
    return false
}

func GetUserByID(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    c.JSON(200, user)
}