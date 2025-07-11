package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

type MemCourse struct {
	ID             uint64
	PassengerID    string
	PassengerName  string
	PassengerEmail string
	StartLat       float64
	StartLng       float64
	EndLat         float64
	EndLng         float64
	DistanceKm     float64
	CO2            float64
	Amount         float64
	Status         string
	DriverID       string
	AcceptedAt     *time.Time
	CreatedAt      time.Time
}

var (
	memCourses []MemCourse
	nextID     uint64
	mu         sync.Mutex
)

func resetTestData() {
	mu.Lock()
	defer mu.Unlock()
	memCourses = nil
	nextID = 1
}

func saveCourse(course *MemCourse) {
	mu.Lock()
	defer mu.Unlock()
	if course.ID == 0 {
		course.ID = nextID
		nextID++
		course.CreatedAt = time.Now()
	}
	found := false
	for i := range memCourses {
		if memCourses[i].ID == course.ID {
			memCourses[i] = *course
			found = true
			break
		}
	}
	if !found {
		memCourses = append(memCourses, *course)
	}
}

func findCourseByID(id uint64) (*MemCourse, bool) {
	mu.Lock()
	defer mu.Unlock()
	for i := range memCourses {
		if memCourses[i].ID == id {
			return &memCourses[i], true
		}
	}
	return nil, false
}

func CreateCourseTest(c *gin.Context) {
	var input struct {
		PassengerID    string  `json:"passenger_id"`
		PassengerName  string  `json:"passenger_name"`
		PassengerEmail string  `json:"passenger_email"`
		StartLat       float64 `json:"start_lat"`
		StartLng       float64 `json:"start_lng"`
		EndLat         float64 `json:"end_lat"`
		EndLng         float64 `json:"end_lng"`
		DistanceKm     float64 `json:"distance_km"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	amount := 4.0 + 0.8*input.DistanceKm
	if !HasEnoughBalanceTest(input.PassengerID, amount) {
		c.JSON(400, gin.H{"error": "Solde insuffisant"})
		return
	}
	course := MemCourse{
		PassengerID:    input.PassengerID,
		PassengerName:  input.PassengerName,
		PassengerEmail: input.PassengerEmail,
		StartLat:       input.StartLat,
		StartLng:       input.StartLng,
		EndLat:         input.EndLat,
		EndLng:         input.EndLng,
		DistanceKm:     input.DistanceKm,
		CO2:            input.DistanceKm * 100,
		Amount:         amount,
		Status:         "requested",
	}
	saveCourse(&course)
	c.JSON(201, course)
}

func GetCourseTest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	course, ok := findCourseByID(id)
	if !ok {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, course)
}

func AcceptCourseTest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input struct {
		DriverID string `json:"driver_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.DriverID == "" {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	course, ok := findCourseByID(id)
	if !ok {
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
	saveCourse(course)
	c.JSON(200, course)
}

func StartCourseTest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	course, ok := findCourseByID(id)
	if !ok {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status != "accepted" {
		c.JSON(400, gin.H{"error": "Course non acceptée ou déjà commencée"})
		return
	}
	course.Status = "in_progress"
	saveCourse(course)
	c.JSON(200, course)
}

func CompleteCourseTest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	course, ok := findCourseByID(id)
	if !ok {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if course.Status != "in_progress" {
		c.JSON(400, gin.H{"error": "Course pas en cours"})
		return
	}
	course.Status = "completed"
	saveCourse(course)
	c.JSON(200, course)
}

func CancelCourseTest(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	course, ok := findCourseByID(id)
	if !ok {
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
	saveCourse(course)
	c.JSON(200, course)
}

func ListCoursesTest(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, memCourses)
}

var HasEnoughBalanceTest = func(googleID string, amount float64) bool { return true }

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/courses", CreateCourseTest)
	r.GET("/courses", ListCoursesTest)
	r.GET("/courses/:id", GetCourseTest)
	r.PATCH("/courses/:id/accept", AcceptCourseTest)
	r.PATCH("/courses/:id/cancel", CancelCourseTest)
	r.PATCH("/courses/:id/start", StartCourseTest)
	r.PATCH("/courses/:id/complete", CompleteCourseTest)
	return r
}

func TestDriverFlow(t *testing.T) {
	resetTestData()
	HasEnoughBalanceTest = func(_ string, _ float64) bool { return true }
	r := setupTestRouter()

	courseData := map[string]interface{}{
		"passenger_id":    "pass1",
		"passenger_name":  "Jean Test",
		"passenger_email": "jean@example.com",
		"start_lat":       48.85,
		"start_lng":       2.35,
		"end_lat":         48.86,
		"end_lng":         2.36,
		"distance_km":     12.4,
	}
	body, _ := json.Marshal(courseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("Course not created, got %d, body: %s", w.Code, w.Body.String())
	}
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	idFloat, ok := created["ID"].(float64)
	if !ok {
		t.Fatalf("No id in response: %v", created)
	}
	id := int(idFloat)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/courses/"+strconv.Itoa(id), nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Course GET failed: %d %s", w.Code, w.Body.String())
	}

	acceptData := map[string]interface{}{"driver_id": "driverX"}
	body, _ = json.Marshal(acceptData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/accept", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Accept failed: %d %s", w.Code, w.Body.String())
	}
	var accepted map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &accepted)
	if accepted["Status"] != "accepted" {
		t.Fatalf("Course not accepted, got status %v", accepted["Status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/start", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Start failed: %d %s", w.Code, w.Body.String())
	}
	var started map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &started)
	if started["Status"] != "in_progress" {
		t.Fatalf("Course not in progress, got %v", started["Status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/complete", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Complete failed: %d %s", w.Code, w.Body.String())
	}
	var completed map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &completed)
	if completed["Status"] != "completed" {
		t.Fatalf("Course not completed, got %v", completed["Status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Fatalf("Cancel should have failed on completed course!")
	}
}

func TestCourseCancel(t *testing.T) {
	resetTestData()
	HasEnoughBalanceTest = func(_ string, _ float64) bool { return true }
	r := setupTestRouter()

	courseData := map[string]interface{}{
		"passenger_id":    "pass1",
		"passenger_name":  "Test",
		"passenger_email": "test@example.com",
		"start_lat":       48.85,
		"start_lng":       2.35,
		"end_lat":         48.86,
		"end_lng":         2.36,
		"distance_km":     10.0,
	}
	body, _ := json.Marshal(courseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	idFloat, ok := created["ID"].(float64)
	if !ok {
		t.Fatalf("No id in response: %v", created)
	}
	id := int(idFloat)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Cancel failed: %d %s", w.Code, w.Body.String())
	}
	var canceled map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &canceled)
	if canceled["Status"] != "cancelled" {
		t.Fatalf("Course not cancelled, got %v", canceled["Status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Fatalf("Should not cancel twice!")
	}
}

func TestInsufficientBalance(t *testing.T) {
	resetTestData()
	HasEnoughBalanceTest = func(string, float64) bool { return false }
	r := setupTestRouter()

	courseData := map[string]interface{}{
		"passenger_id":    "nopass",
		"passenger_name":  "Pauvre",
		"passenger_email": "pauvre@example.com",
		"start_lat":       0.0,
		"start_lng":       0.0,
		"end_lat":         1.0,
		"end_lng":         1.0,
		"distance_km":     2.0,
	}
	body, _ := json.Marshal(courseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Course creation should fail for insufficient balance")
	}
}
