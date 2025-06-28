package main

import (
    "testing"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "strings"
)

func generateJWT(roles []string) string {
    claims := jwt.MapClaims{
        "sub":   "testid",
        "roles": roles,
        "exp":   time.Now().Add(time.Hour).Unix(),
        "email": "liam2106cariou@gmail.com",
    }
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := jwtToken.SignedString([]byte("devsecret")) 
    return tokenString
}

func TestAuthAdminSecurity(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.Use(AuthRequired())
    admin := r.Group("/admin")
    admin.Use(AdminRequired())
    admin.GET("/users", func(c *gin.Context) { c.JSON(200, gin.H{"ok":true}) })

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/admin/users", nil)
    r.ServeHTTP(w, req)
    if w.Code != 401 && w.Code != 403 {
        t.Fatalf("Expected 401/403 for no token, got %d", w.Code)
    }

    w2 := httptest.NewRecorder()
    req2 := httptest.NewRequest("GET", "/admin/users", nil)
    req2.Header.Set("Authorization", "Bearer "+generateJWT([]string{"ROLE_PASSAGER"}))
    r.ServeHTTP(w2, req2)
    if w2.Code != 403 {
        t.Fatalf("Expected 403 for not admin, got %d", w2.Code)
    }

    w3 := httptest.NewRecorder()
    req3 := httptest.NewRequest("GET", "/admin/users", nil)
    req3.Header.Set("Authorization", "Bearer "+generateJWT([]string{"ROLE_ADMIN"}))
    r.ServeHTTP(w3, req3)
    if w3.Code != 200 {
        t.Fatalf("Expected 200 for admin, got %d", w3.Code)
    }
}