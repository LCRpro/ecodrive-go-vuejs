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
	driverView := c.Query("driverView")
	driverID := c.Query("driver_id")
	status := c.Query("status")
	var courses []Course

	if driverView == "1" && driverID != "" {
		db.Where("driver_id = ? AND (status = ? OR status = ?)", driverID, "accepted", "in_progress").
			Order("created_at desc").Find(&courses)
		if len(courses) > 0 {
			c.JSON(200, courses)
			return
		}
		db.Where("status = ?", "requested").Order("created_at desc").Find(&courses)
		c.JSON(200, courses)
		return
	} else if status != "" {
		db.Where("status = ?", status).Order("created_at desc").Find(&courses)
	} else if role == "passenger" && id != "" {
		db.Where("passenger_id = ?", id).Order("created_at desc").Find(&courses)
	} else if role == "driver" && id != "" {
		db.Where("driver_id = ?", id).Order("created_at desc").Find(&courses)
	} else {
		db.Order("created_at desc").Find(&courses)
	}
	c.JSON(200, courses)
}

func GetCourse(c *gin.Context) {
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	var course Course
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
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	var course Course
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
	if err := db.Save(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
	c.JSON(200, course)
}

func CancelCourse(c *gin.Context) {
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	var course Course
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status == "cancelled" {
		c.JSON(400, gin.H{"error": "Déjà annulée"})
		return
	}
	if course.Status == "completed" {
		c.JSON(400, gin.H{"error": "Course déjà terminée, annulation impossible"})
		return
	}
	course.Status = "cancelled"
	if err := db.Save(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
	c.JSON(200, course)
}

var HasEnoughBalance = func(googleID string, amount float64) bool {
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
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	var course Course
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status != "in_progress" {
		c.JSON(400, gin.H{"error": "Course pas en cours"})
		return
	}
	course.Status = "completed"
	if err := db.Save(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
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
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	var course Course
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
	if err := db.Save(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
	c.JSON(200, course)
}

func PatchCourse(c *gin.Context) {
	var update map[string]interface{}
	idU64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID invalide"})
		return
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	var course Course
	if err := db.Where("id = ?", idU64).First(&course).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	for k, v := range update {
		switch k {
		case "passenger_id":
			course.PassengerID = v.(string)
		case "passenger_name":
			course.PassengerName = v.(string)
		case "passenger_email":
			course.PassengerEmail = v.(string)
		case "start_lat":
			course.StartLat = v.(float64)
		case "start_lng":
			course.StartLng = v.(float64)
		case "end_lat":
			course.EndLat = v.(float64)
		case "end_lng":
			course.EndLng = v.(float64)
		case "distance_km":
			course.DistanceKm = v.(float64)
		case "co2":
			course.CO2 = v.(float64)
		case "amount":
			course.Amount = v.(float64)
		case "status":
			course.Status = v.(string)
		case "driver_id":
			course.DriverID = v.(string)
		}
	}
	if err := db.Save(&course).Error; err != nil {
		c.JSON(500, gin.H{"error": "save failed"})
		return
	}
	c.JSON(200, course)
}