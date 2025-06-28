package main

import (
    "github.com/gin-gonic/gin"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "os"
    "github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

var (
jwtSecret = []byte(os.Getenv("JWT_SECRET"))
    googleOauthConfig *oauth2.Config
)

func main() {
    _ = godotenv.Load()
    googleOauthConfig = &oauth2.Config{
        RedirectURL:  "http://localhost:8000/auth/callback",
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
        Endpoint:     google.Endpoint,
    }

    r := gin.Default()
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
    r.GET("/auth/google", HandleGoogleAuth)
    r.GET("/auth/callback", HandleGoogleCallback)
    r.Run(":8000")
}