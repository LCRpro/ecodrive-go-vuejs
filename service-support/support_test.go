package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func resetTicketsData() {
	tickets = nil
	nextTicketID = 1
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/support", CreateSupportTicket)
	r.GET("/support/user/:google_id", ListUserTickets)
	r.GET("/support/all", ListAllTickets)
	r.PATCH("/support/reply/:id", ReplySupportTicket)
	return r
}

func TestCreateSupportTicket(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	ticket := SupportTicket{
		GoogleID: "test-google-id",
		Category: "bug",
		Message:  "Test message",
	}
	jsonValue, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("POST", "/support", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", w.Code)
	}
}

func TestCreateSupportTicket_BadJSON(t *testing.T) {
	resetTicketsData()
	router := setupRouter()
	req, _ := http.NewRequest("POST", "/support", bytes.NewBuffer([]byte("{badjson")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestCreateSupportTicket_MissingFields(t *testing.T) {
	resetTicketsData()
	router := setupRouter()
	ticket := SupportTicket{
		GoogleID: "",
		Category: "",
		Message:  "",
	}
	jsonValue, _ := json.Marshal(ticket)
	req, _ := http.NewRequest("POST", "/support", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListUserTickets(t *testing.T) {
	resetTicketsData()
	ticket := SupportTicket{
		GoogleID:  "user1",
		Category:  "help",
		Message:   "message",
		Status:    "ouvert",
		CreatedAt: time.Now(),
	}
	saveTicket(&ticket)

	gin.SetMode(gin.TestMode)
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/support/user/user1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
}

func TestListUserTickets_UserNotFound(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/support/user/inconnu", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
	var result []SupportTicket
	json.Unmarshal(w.Body.Bytes(), &result)
	if len(result) != 0 {
		t.Errorf("Expected empty result for unknown user")
	}
}

func TestListAllTickets(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/support/all", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
}

func TestReplySupportTicket(t *testing.T) {
	resetTicketsData()
	ticket := SupportTicket{
		GoogleID:  "user2",
		Category:  "question",
		Message:   "Need help",
		Status:    "ouvert",
		CreatedAt: time.Now(),
	}
	saveTicket(&ticket)
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	reply := map[string]string{
		"admin_reply": "We fixed it",
		"status":      "fermé",
	}
	jsonValue, _ := json.Marshal(reply)
	url := "/support/reply/" + strconv.Itoa(int(ticket.ID))
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
}

func TestReplySupportTicket_BadID(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	reply := map[string]string{
		"admin_reply": "test",
		"status":      "ouvert",
	}
	jsonValue, _ := json.Marshal(reply)
	req, _ := http.NewRequest("PATCH", "/support/reply/abc", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest, got %d", w.Code)
	}
}

func TestReplySupportTicket_NotFound(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	reply := map[string]string{
		"admin_reply": "test",
		"status":      "ouvert",
	}
	jsonValue, _ := json.Marshal(reply)
	req, _ := http.NewRequest("PATCH", "/support/reply/9999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 NotFound, got %d", w.Code)
	}
}

func TestReplySupportTicket_BadJSON(t *testing.T) {
	resetTicketsData()
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	req, _ := http.NewRequest("PATCH", "/support/reply/1", bytes.NewBuffer([]byte("{badjson")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest, got %d", w.Code)
	}
}

func TestReplySupportTicket_NoReply(t *testing.T) {
	resetTicketsData()
	ticket := SupportTicket{
		GoogleID:  "user2",
		Category:  "question",
		Message:   "Need help",
		Status:    "ouvert",
		CreatedAt: time.Now(),
	}
	saveTicket(&ticket)
	gin.SetMode(gin.TestMode)
	router := setupRouter()
	reply := map[string]string{
		"admin_reply": "",
		"status":      "ouvert",
	}
	jsonValue, _ := json.Marshal(reply)
	url := "/support/reply/" + strconv.Itoa(int(ticket.ID))
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest, got %d", w.Code)
	}
}
