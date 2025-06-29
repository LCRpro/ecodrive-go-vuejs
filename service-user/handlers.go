package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "fmt"
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

func ListUsers(c *gin.Context) {
    var users []User
    db.Find(&users)
    c.JSON(200, users)
}

func PatchUser(c *gin.Context) {
    var user User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    var input map[string]interface{}
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    db.Model(&user).Updates(input)
    db.First(&user, user.ID) 
    c.JSON(200, user)
}

func CreateDriverRequest(c *gin.Context) {
    var input struct {
        GoogleID string `json:"google_id"`
        Car      string `json:"car"`
        Plate    string `json:"plate"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }

    var user User
    if err := db.Where("google_id = ?", input.GoogleID).First(&user).Error; err == nil {
        db.Model(&user).Updates(map[string]interface{}{"car": input.Car, "plate": input.Plate})
    }

    var existing DriverRequest
    if err := db.Where("google_id = ? AND status = ?", input.GoogleID, "pending").First(&existing).Error; err == nil {
        c.JSON(400, gin.H{"error": "Déjà une demande en attente"})
        return
    }

    req := DriverRequest{
        GoogleID: input.GoogleID,
        Car:      input.Car,
        Plate:    input.Plate,
        Status:   "pending",
    }
    db.Create(&req)
    c.JSON(201, req)
}

func ListDriverRequests(c *gin.Context) {
    var requests []DriverRequest
    db.Order("status desc, id asc").Find(&requests)
    c.JSON(200, requests)
}

func HandleDriverRequest(c *gin.Context) {
    id := c.Param("id")
    var req DriverRequest
    if err := db.First(&req, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    var input struct{ Action string }
    if err := c.ShouldBindJSON(&input); err != nil || (input.Action != "accept" && input.Action != "refuse") {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    if req.Status != "pending" {
        c.JSON(400, gin.H{"error": "Déjà traité"})
        return
    }

    if input.Action == "accept" {
        var user User
        if err := db.Where("google_id = ?", req.GoogleID).First(&user).Error; err == nil {
            var roles []string
            json.Unmarshal([]byte(user.Roles), &roles)
            if !contains(roles, "ROLE_DRIVER") {
                roles = append(roles, "ROLE_DRIVER")
                rolesJSON, _ := json.Marshal(roles)
                user.Roles = string(rolesJSON)
            }
            user.Car = req.Car
            user.Plate = req.Plate
            db.Save(&user)
        }
        req.Status = "accepted"
        fmt.Printf("Accept driver for %s (%s)\n", user.Email, user.GoogleID)
    } else {
        req.Status = "refused"
    }
    
    db.Save(&req)
    c.JSON(200, req)
}
