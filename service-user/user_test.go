package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

var (
	users      []User
	nextUserID uint = 1
)

func resetTestData() {
	users = nil
	nextUserID = 1
}

func FindOrCreateUserTest(c *gin.Context) {
	var in struct {
		GoogleID string `json:"google_id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	for i := range users {
		if users[i].Email == in.Email {
			users[i].GoogleID = in.GoogleID
			c.JSON(200, users[i])
			return
		}
	}
	for i := range users {
		if users[i].Email == in.Email {
			users[i].GoogleID = in.GoogleID
			c.JSON(200, users[i])
			return
		}
	}
	roles := `["ROLE_PASSAGER"]`
	if in.Email == "liam2106cariou@gmail.com" || in.Name == "Liam Cariou" {
		roles = `["ROLE_ADMIN","ROLE_PASSAGER"]`
	}
	user := User{
		ID:       nextUserID,
		GoogleID: in.GoogleID,
		Email:    in.Email,
		Name:     in.Name,
		Roles:    roles,
	}
	nextUserID++
	users = append(users, user)
	c.JSON(200, user)
}

func ListUsersTest(c *gin.Context) {
	c.JSON(200, users)
}

func PatchUserTest(c *gin.Context) {
	gid := c.Param("id")
	for i := range users {
		if users[i].GoogleID == gid {
			var input map[string]interface{}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(400, gin.H{"error": "bad input"})
				return
			}
			if v, ok := input["balance"]; ok {
				if f, ok := v.(float64); ok {
					users[i].Balance = f
				}
			}
			c.JSON(200, users[i])
			return
		}
	}
	c.JSON(404, gin.H{"error": "not found"})
}

func GetUserByIDTest(c *gin.Context) {
	gid := c.Param("id")
	for _, u := range users {
		if u.GoogleID == gid {
			c.JSON(200, u)
			return
		}
	}
	c.JSON(404, gin.H{"error": "not found"})
}

func setupRouterUserTest() *gin.Engine {
	r := gin.New()
	r.POST("/users/find-or-create", FindOrCreateUserTest)
	r.GET("/users", ListUsersTest)
	r.PATCH("/users/:id", PatchUserTest)
	r.GET("/users/:id", GetUserByIDTest)
	return r
}

func TestFindOrCreateUser(t *testing.T) {
	resetTestData()
	r := setupRouterUserTest()
	body1 := `{"google_id":"","email":"mail@test.com","name":"Toto"}`
	req1 := httptest.NewRequest("POST", "/users/find-or-create", bytes.NewBufferString(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != 200 {
		t.Fatalf("Create user fail: %d %s", w1.Code, w1.Body.String())
	}

	body2 := `{"google_id":"gid1","email":"mail@test.com","name":"Toto"}`
	req2 := httptest.NewRequest("POST", "/users/find-or-create", bytes.NewBufferString(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != 200 {
		t.Fatalf("FindOrCreate (update google_id on email) fail: %d", w2.Code)
	}

	req3 := httptest.NewRequest("POST", "/users/find-or-create", bytes.NewBufferString(body2))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != 200 {
		t.Fatalf("FindOrCreate (exist) fail: %d", w3.Code)
	}
}

func TestListUsers(t *testing.T) {
	resetTestData()
	users = append(users, User{ID: 1, GoogleID: "a", Email: "a", Name: "A"})
	r := setupRouterUserTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("ListUsers fail: %d", w.Code)
	}
}

func TestPatchUser(t *testing.T) {
	resetTestData()
	users = append(users, User{ID: 1, GoogleID: "patchid", Name: "Pat", Email: "pat@test.com", Balance: 1.5})
	r := setupRouterUserTest()
	patch := `{"balance":42}`
	req := httptest.NewRequest("PATCH", "/users/patchid", bytes.NewBufferString(patch))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("PatchUser fail: %d", w.Code)
	}
	req2 := httptest.NewRequest("PATCH", "/users/unknown", bytes.NewBufferString(patch))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != 404 {
		t.Fatalf("PatchUser (not found) should return 404, got %d", w2.Code)
	}
}

func TestGetUserByID(t *testing.T) {
	resetTestData()
	users = append(users, User{ID: 1, GoogleID: "userget", Name: "Get", Email: "get@get"})
	r := setupRouterUserTest()
	req := httptest.NewRequest("GET", "/users/userget", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("GetUserByID fail: %d", w.Code)
	}
	req2 := httptest.NewRequest("GET", "/users/notfound", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != 404 {
		t.Fatalf("GetUserByID (not found) should return 404, got %d", w2.Code)
	}
}
