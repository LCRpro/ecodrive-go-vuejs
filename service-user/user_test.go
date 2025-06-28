package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "net/http/httptest"
    "strings"
    "testing"
)

var testdb *gorm.DB

func setupRouter() *gin.Engine {
    r := gin.Default()

    r.POST("/users", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(400, gin.H{"error": "invalid body"})
            return
        }
        testdb.Create(&user)
        c.JSON(201, user)
    })

    r.GET("/users/:id", func(c *gin.Context) {
        var user User
        if err := testdb.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(404, gin.H{"error": "not found"})
            return
        }
        c.JSON(200, user)
    })

    return r
}

func TestUserAPI(t *testing.T) {
    var err error
    testdb, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open test db: %v", err)
    }
    testdb.AutoMigrate(&User{})

    router := setupRouter()

    w := httptest.NewRecorder()
    req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"GoogleID":"g123","Email":"test@tdd.com","Name":"Jean Test"}`))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    if w.Code != 201 {
        t.Fatalf("expected 201, got %d", w.Code)
    }

    var createdUser User
    testdb.First(&createdUser)

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest("GET", "/users/"+string(rune(createdUser.ID)), nil)
    router.ServeHTTP(w2, req2)
    if w2.Code != 200 {
        t.Fatalf("expected 200, got %d", w2.Code)
    }
}