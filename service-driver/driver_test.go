package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/courses", CreateCourse)
	r.GET("/courses", ListCourses)
	r.GET("/courses/:id", GetCourse)
	r.PATCH("/courses/:id/accept", AcceptCourse)
	r.PATCH("/courses/:id/cancel", CancelCourse)
	r.PATCH("/courses/:id/start", StartCourse)
	r.PATCH("/courses/:id/complete", CompleteCourse)
	return r
}

func setupTestDB() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test db")
	}
	db.AutoMigrate(&Course{})
}

func mockBalanceOK(googleID string, amount float64) bool { return true }

func TestDriverFlow(t *testing.T) {
	setupTestDB()
	HasEnoughBalance = mockBalanceOK

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
		"co2":             0.0,
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
	id := int(created["id"].(float64))

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
	if accepted["status"] != "accepted" {
		t.Fatalf("Course not accepted, got status %v", accepted["status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/start", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Start failed: %d %s", w.Code, w.Body.String())
	}
	var started map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &started)
	if started["status"] != "in_progress" {
		t.Fatalf("Course not in progress, got %v", started["status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/complete", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Complete failed: %d %s", w.Code, w.Body.String())
	}
	var completed map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &completed)
	if completed["status"] != "completed" {
		t.Fatalf("Course not completed, got %v", completed["status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Fatalf("Cancel should have failed on completed course!")
	}
}

func TestCourseCancel(t *testing.T) {
	setupTestDB()
	HasEnoughBalance = mockBalanceOK
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
		"co2":             0.0,
	}
	body, _ := json.Marshal(courseData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/courses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	id := int(created["id"].(float64))

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Cancel failed: %d %s", w.Code, w.Body.String())
	}
	var canceled map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &canceled)
	if canceled["status"] != "cancelled" {
		t.Fatalf("Course not cancelled, got %v", canceled["status"])
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/courses/"+strconv.Itoa(id)+"/cancel", nil)
	r.ServeHTTP(w, req)
	if w.Code == 200 {
		t.Fatalf("Should not cancel twice!")
	}
}

func TestInsufficientBalance(t *testing.T) {
	setupTestDB()
	HasEnoughBalance = func(string, float64) bool { return false }
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
		"co2":             0.0,
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