package main

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "os"
    "strings"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Bearer token"})
            return
        }
        tokenString := strings.TrimPrefix(auth, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
            return
        }
        c.Set("claims", claims)
        c.Next()
    }
}

func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        claimsVal, exists := c.Get("claims")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing claims"})
            return
        }
        claims := claimsVal.(jwt.MapClaims)
        var roles []string
        raw := claims["roles"]
        switch v := raw.(type) {
        case []interface{}:
            for _, r := range v { roles = append(roles, r.(string)) }
        case []string:
            roles = v
        }
        isAdmin := false
        for _, r := range roles {
            if r == "ROLE_ADMIN" {
                isAdmin = true
                break
            }
        }
        if !isAdmin {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin only"})
            return
        }
        c.Next()
    }
}