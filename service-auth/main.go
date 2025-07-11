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
	jwtSecret           []byte
	googleOauthConfig   *oauth2.Config
	userServiceURL      string
	frontendCallbackURL string
)

func main() {
	_ = godotenv.Load()
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	userServiceURL = os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8002"
	}

	frontendCallbackURL = os.Getenv("FRONTEND_CALLBACK_URL")
	if frontendCallbackURL == "" {
		frontendCallbackURL = "http://localhost:5173/callback"
	}

	authRedirectURL := os.Getenv("AUTH_REDIRECT_URL")
	if authRedirectURL == "" {
		authRedirectURL = "http://localhost:8000/auth/callback"
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  authRedirectURL,
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:5173"
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.GET("/auth/google", HandleGoogleAuth)
	r.GET("/auth/callback", HandleGoogleCallback)
	r.Run(":8000")
}
