package main

import (
    "encoding/json"
    "github.com/gin-gonic/gin"
    "fmt"
)

func contains(list []string, item string) bool {
    for _, v := range list {
        if v == item {
            return true
        }
    }
    return false
}

func FindOrCreateUser(c *gin.Context) {
    var in struct {
        GoogleID string `json:"google_id"`
        Email    string `json:"email"`
        Name     string `json:"name"`
    }
    if err := c.ShouldBindJSON(&in); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }

    var user User
    if err := db.Where("google_id = ?", in.GoogleID).First(&user).Error; err == nil {
        c.JSON(200, user)
        return
    }
    if err := db.Where("email = ?", in.Email).First(&user).Error; err == nil {
        user.GoogleID = in.GoogleID
        db.Save(&user)
        c.JSON(200, user)
        return
    }

    var roles string
    if in.Email == "liam2106cariou@gmail.com" || in.Name == "Liam Cariou" {
        roles = `["ROLE_ADMIN","ROLE_PASSAGER"]`
    } else {
        roles = `["ROLE_PASSAGER"]`
    }

    user = User{
        GoogleID: in.GoogleID,
        Email:    in.Email,
        Name:     in.Name,
        Roles:    roles,
    }
    db.Create(&user)
    c.JSON(200, user)
}

func ListUsers(c *gin.Context) {
    var users []User
    db.Find(&users)
    c.JSON(200, users)
}

func PatchUser(c *gin.Context) {
    var user User
    if err := db.Where("google_id = ?", c.Param("id")).First(&user).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    var input map[string]interface{}
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    if v, ok := input["balance"]; ok {
        if f, ok := v.(float64); ok {
            user.Balance = f
        }
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
func GetUserByID(c *gin.Context) {
    var user User
    if err := db.Where("google_id = ?", c.Param("id")).First(&user).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    c.JSON(200, user)
}

func GetAppBalance(c *gin.Context) {
    var acc AppAccount
    if err := db.First(&acc).Error; err != nil {
        c.JSON(500, gin.H{"error": "no app account"})
        return
    }
    c.JSON(200, acc)
}

func CreditAppAccount(c *gin.Context) {
    var input struct {
        Amount float64 `json:"amount"`
    }
    if err := c.ShouldBindJSON(&input); err != nil || input.Amount == 0 {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }

    var app AppAccount
    if err := db.First(&app, 1).Error; err != nil {
        app = AppAccount{ID: 1, Balance: input.Amount}
        db.Create(&app)
    } else {
        app.Balance += input.Amount
        db.Save(&app)
    }
    c.JSON(200, app)
}


func PatchDriverRequest(c *gin.Context) {
    id := c.Param("id")
    var req DriverRequest
    if err := db.First(&req, id).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    var input map[string]interface{}
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    db.Model(&req).Updates(input)
    db.First(&req, id)
    c.JSON(200, req)
}