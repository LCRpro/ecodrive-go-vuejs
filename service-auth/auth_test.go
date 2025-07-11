package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var testJWTSecret = []byte("test-secret")

func setupTestGoogleOauthConfig() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     "fake-client",
		ClientSecret: "fake-secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: frontendCallbackURL,
		Scopes:      []string{"email", "profile"},
	}
}

type fakeRoundTripper struct {
	handler func(req *http.Request) *http.Response
}

func (f *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.handler(req), nil
}

func newMockClient(handler func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &fakeRoundTripper{handler: handler},
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/auth/google", HandleGoogleAuth)
	r.GET("/auth/google/callback", HandleGoogleCallback)
	return r
}

func TestJWTGeneration(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":   "test-id",
		"email": "test@tdd.com",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(testJWTSecret)
	if err != nil {
		t.Fatalf("Erreur génération JWT: %v", err)
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testJWTSecret, nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("Le JWT n'est pas valide")
	}
	claims2, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims2["email"] != "test@tdd.com" {
		t.Fatal("Claims incorrects dans le JWT")
	}
}

func TestHandleGoogleAuth_Redirects(t *testing.T) {
	setupTestGoogleOauthConfig()
	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/google", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307 redirect, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "accounts.google.com") {
		t.Errorf("Redirection OAuth attendue vers Google, obtenu : %s", loc)
	}
}

func TestHandleGoogleCallback_Success(t *testing.T) {
	setupTestGoogleOauthConfig()
	http.DefaultClient = newMockClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "oauth2.googleapis.com/token") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-access","token_type":"Bearer","expires_in":3600}`)),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			user := GoogleUser{Id: "gid", Email: "mail@x.com", Name: "toto"}
			body, _ := json.Marshal(user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "/users/find-or-create") {
			user := User{GoogleID: "gid", Email: "mail@x.com", Name: "toto", Roles: `["user","admin"]`}
			body, _ := json.Marshal(user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}
	})
	jwtSecret = testJWTSecret

	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/google/callback?code=fake", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Expected 307, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, frontendCallbackURL+"?token=") {
		t.Errorf("Redirection inattendue : %s", loc)
	}

	parts := strings.Split(loc, "token=")
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testJWTSecret, nil
	})
	if err != nil || !token.Valid {
		t.Fatal("JWT invalide dans callback")
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	if claims["email"] != "mail@x.com" {
		t.Errorf("Email manquant dans le JWT")
	}
	if claims["sub"] != "gid" {
		t.Errorf("sub incorrect dans le JWT")
	}
}

func TestHandleGoogleCallback_FailUserInfo(t *testing.T) {
	setupTestGoogleOauthConfig()
	http.DefaultClient = newMockClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "oauth2.googleapis.com/token") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-access","token_type":"Bearer","expires_in":3600}`)),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("fail userinfo")),
				Header:     make(http.Header),
			}
		}
		return nil
	})
	jwtSecret = testJWTSecret
	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/google/callback?code=fail", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Attendait 500, obtenu %d", w.Code)
	}
}

func TestHandleGoogleCallback_FailUserService(t *testing.T) {
	setupTestGoogleOauthConfig()
	http.DefaultClient = newMockClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "oauth2.googleapis.com/token") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-access","token_type":"Bearer","expires_in":3600}`)),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			user := GoogleUser{Id: "gid", Email: "mail@x.com", Name: "toto"}
			body, _ := json.Marshal(user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "/users/find-or-create") {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("fail user service")),
				Header:     make(http.Header),
			}
		}
		return nil
	})
	jwtSecret = testJWTSecret
	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/google/callback?code=x", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Attendait 500, obtenu %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Erreur service-user") {
		t.Errorf("Erreur user service manquante")
	}
}

func TestHandleGoogleCallback_BadRolesJSON(t *testing.T) {
	setupTestGoogleOauthConfig()
	http.DefaultClient = newMockClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.String(), "oauth2.googleapis.com/token") {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"fake-access","token_type":"Bearer","expires_in":3600}`)),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "userinfo") {
			user := GoogleUser{Id: "gid", Email: "mail@x.com", Name: "toto"}
			body, _ := json.Marshal(user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		if strings.Contains(req.URL.String(), "/users/find-or-create") {
			user := User{
				GoogleID: "gid", Email: "mail@x.com", Name: "toto",
				Roles: `invalid-json`,
			}
			body, _ := json.Marshal(user)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(string(body))),
				Header:     make(http.Header),
			}
		}
		return nil
	})
	jwtSecret = testJWTSecret
	router := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/auth/google/callback?code=fake", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Expected 307, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.HasPrefix(loc, frontendCallbackURL+"?token=") {
		t.Errorf("Redirection inattendue : %s", loc)
	}
	parts := strings.Split(loc, "token=")
	tokenString := parts[1]
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return testJWTSecret, nil
	})
	claims, _ := token.Claims.(jwt.MapClaims)
	if roles, ok := claims["roles"].([]interface{}); !ok || len(roles) != 0 {
		t.Errorf("roles devrait être vide si JSON mal formé")
	}
}
