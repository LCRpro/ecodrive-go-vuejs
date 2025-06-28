package main

import (
    "testing"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)

func TestJWTGeneration(t *testing.T) {
    claims := jwt.MapClaims{
        "sub": "test-id",
        "email": "test@tdd.com",
        "exp": time.Now().Add(time.Hour * 24).Unix(),
    }
    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := jwtToken.SignedString(jwtSecret)
    if err != nil {
        t.Fatalf("Erreur génération JWT: %v", err)
    }
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        t.Fatalf("Le JWT n'est pas valide")
    }

	package main


func TestHandleGoogleAuth(t *testing.T) {
    router := gin.Default()
    router.GET("/auth/google", HandleGoogleAuth)
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/auth/google", nil)
    router.ServeHTTP(w, req)
    if w.Code != 307 {
        t.Fatalf("expected 307 redirect, got %d", w.Code)
    }
}
}