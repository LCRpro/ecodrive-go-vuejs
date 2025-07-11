package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	jwtSecret         []byte // juste déclaré ici
	googleOauthConfig *oauth2.Config
)

func main() {
	_ = godotenv.Load()                         // charge bien .env
	jwtSecret = []byte(os.Getenv("JWT_SECRET")) // ICI on lit la vraie valeur de .env

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "https://auth-ecodrive.liamcariou.fr/auth/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ecodrive.liamcariou.fr"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.GET("/auth/google", HandleGoogleAuth)
	r.GET("/auth/callback", HandleGoogleCallback)
	r.Run(":8000")
}
