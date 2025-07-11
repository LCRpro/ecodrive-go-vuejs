package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type GoogleUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type User struct {
	ID       uint   `json:"id"`
	GoogleID string `json:"google_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Roles    string `json:"roles"`
}

func HandleGoogleAuth(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed exchange: %s", err)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed userinfo: %s", err)
		return
	}
	defer resp.Body.Close()
	var gu GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil {
		c.String(http.StatusInternalServerError, "Failed decode: %s", err)
		return
	}

	user, err := findOrCreateUser(gu)
	if err != nil {
		c.String(http.StatusInternalServerError, "Erreur service-user: %s", err)
		return
	}

	var roles []string
	if err := json.Unmarshal([]byte(user.Roles), &roles); err != nil {
		roles = []string{}
	}

	claims := jwt.MapClaims{
		"sub":   user.GoogleID,
		"roles": roles,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"email": user.Email,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtSecret)
	if err != nil {
		c.String(http.StatusInternalServerError, "JWT error: %s", err)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, frontendCallbackURL+"?token="+tokenString)
}

func findOrCreateUser(gu GoogleUser) (User, error) {
	payload, _ := json.Marshal(gin.H{
		"google_id": gu.Id,
		"email":     gu.Email,
		"name":      gu.Name,
	})
	resp, err := http.Post(userServiceURL+"/users/find-or-create", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return User{}, fmt.Errorf("service-user status %d", resp.StatusCode)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
