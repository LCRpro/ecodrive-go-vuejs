package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
	"strconv"
)

func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
    if err != nil {
        panic("failed to connect test db")
    }
    db.AutoMigrate(&SupportTicket{})
    return db
}

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.POST("/support", CreateSupportTicket)
    r.GET("/support/user/:google_id", ListUserTickets)
    r.GET("/support/all", ListAllTickets)
    r.PATCH("/support/reply/:id", ReplySupportTicket)
    return r
}

func TestMain(m *testing.M) {
    db = setupTestDB() 
    os.Exit(m.Run())
}

func TestCreateSupportTicket(t *testing.T) {
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

func TestListUserTickets(t *testing.T) {
    router := setupRouter()

    db.Create(&SupportTicket{
        GoogleID: "user1",
        Category: "help",
        Message:  "message",
        Status:   "ouvert",
    })

    req, _ := http.NewRequest("GET", "/support/user/user1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Code)
    }
}

func TestListAllTickets(t *testing.T) {
    router := setupRouter()

    req, _ := http.NewRequest("GET", "/support/all", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected 200 OK, got %d", w.Code)
    }
}

func TestReplySupportTicket(t *testing.T) {
    router := setupRouter()

    ticket := SupportTicket{
        GoogleID: "user2",
        Category: "question",
        Message:  "Need help",
        Status:   "ouvert",
    }
    db.Create(&ticket)

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