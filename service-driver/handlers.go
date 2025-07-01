package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"encoding/json"
	"net/http"
	"strconv"
	"bytes"
	"fmt"
)

func CreateCourse(c *gin.Context) {
	type CourseInput struct {
		PassengerID    string  `json:"passenger_id"`
		PassengerName  string  `json:"passenger_name"`
		PassengerEmail string  `json:"passenger_email"`
		StartLat       float64 `json:"start_lat"`
		StartLng       float64 `json:"start_lng"`
		EndLat         float64 `json:"end_lat"`
		EndLng         float64 `json:"end_lng"`
		DistanceKm     float64 `json:"distance_km"`
	}

	var input CourseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}

co2Saved := input.DistanceKm * 100 

amount := 4.0 + 0.8*input.DistanceKm

if !HasEnoughBalance(input.PassengerID, amount) {
    c.JSON(400, gin.H{"error": "Solde insuffisant pour commander cette course"})
    return
}

course := Course{
    PassengerID:    input.PassengerID,
    PassengerName:  input.PassengerName,
    PassengerEmail: input.PassengerEmail,
    StartLat:       input.StartLat,
    StartLng:       input.StartLng,
    EndLat:         input.EndLat,
    EndLng:         input.EndLng,
    DistanceKm:     input.DistanceKm,
    CO2:            co2Saved,
    Amount:         amount,
    Status:         "requested",
}

	if err := db.Create(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "Erreur lors de la création"})
		return
	}

	c.JSON(201, course)
}

func ListCourses(c *gin.Context) {
	role := c.Query("role")
	id := c.Query("id")
	var courses []Course
	if role == "passenger" && id != "" {
		db.Where("passenger_id = ?", id).Order("created_at desc").Find(&courses)
	} else if role == "driver" && id != "" {
		db.Where("status = ? OR driver_id = ?", "requested", id).Order("created_at desc").Find(&courses)
	} else {
		db.Order("created_at desc").Find(&courses)
	}
	c.JSON(200, courses)
}

func GetCourse(c *gin.Context) {
	var course Course
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, course)
}

func AcceptCourse(c *gin.Context) {
	var input struct {
		DriverID string `json:"driver_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.DriverID == "" {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	var course Course
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status != "requested" {
		c.JSON(400, gin.H{"error": "Course déjà acceptée"})
		return
	}
	now := time.Now()
	course.Status = "accepted"
	course.DriverID = input.DriverID
	course.AcceptedAt = &now
	db.Save(&course)
	c.JSON(200, course)
}

func CancelCourse(c *gin.Context) {
	var course Course
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status == "cancelled" {
		c.JSON(400, gin.H{"error": "Déjà annulée"})
		return
	}
	course.Status = "cancelled"
	db.Save(&course)
	c.JSON(200, course)
}

func HasEnoughBalance(googleID string, amount float64) bool {
	resp, err := http.Get("http://localhost:8002/users/" + googleID)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var user struct{ Balance float64 }
	json.NewDecoder(resp.Body).Decode(&user)
	return user.Balance >= amount
}


func UpdateUserBalance(googleID string, delta float64) {
    resp, err := http.Get("http://localhost:8002/users/" + googleID)
    if err != nil { return }
    defer resp.Body.Close()
    var user struct{ Balance float64 }
    json.NewDecoder(resp.Body).Decode(&user)
    newBalance := user.Balance + delta
    body, _ := json.Marshal(map[string]float64{"balance": newBalance})
    req, _ := http.NewRequest("PATCH", "http://localhost:8002/users/"+googleID, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    client.Do(req)
}

func CompleteCourse(c *gin.Context) {
    var course Course
    idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "ID invalide"})
        return
    }
    if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    if course.Status != "in_progress" {
        c.JSON(400, gin.H{"error": "Course pas en cours"})
        return
    }
    course.Status = "completed"
    db.Save(&course)
    c.JSON(200, course)
}

func creditAppAccount(amount float64) {
    body := []byte(fmt.Sprintf(`{"amount":%f}`, amount))
    req, _ := http.NewRequest("PATCH", "http://localhost:8002/app-account/credit", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    client.Do(req)
}

func StartCourse(c *gin.Context) {
    var course Course
    idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(400, gin.H{"error": "ID invalide"})
        return
    }
    if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }
    if course.Status != "accepted" {
        c.JSON(400, gin.H{"error": "Course non acceptée ou déjà commencée"})
        return
    }

    go UpdateUserBalance(course.PassengerID, -course.Amount)

if course.DriverID != "" {
    driverGain := course.Amount * 0.80  
    go UpdateUserBalance(course.DriverID, driverGain)
}

appShare := course.Amount * 0.20 
go creditAppAccount(appShare)

    course.Status = "in_progress"
    db.Save(&course)
    c.JSON(200, course)
}