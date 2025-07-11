package main

import (
    "encoding/json"
    "net/http"
    "github.com/gin-gonic/gin"
    "bytes"
)

type User struct {
    ID        uint   `json:"id"`
    GoogleID  string `json:"google_id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    Roles     string `json:"roles"`
    Birthdate string `json:"birthdate"`
    Gender    string `json:"gender"`
    Car       string `json:"car"`
    Plate     string `json:"plate"`
    Phone     string `json:"phone"` 
}

func ListUsers(c *gin.Context) {
    resp, err := http.Get("http://localhost:8002/users")
    if err != nil {
        c.JSON(500, gin.H{"error": "service-user injoignable"})
        return
    }
    defer resp.Body.Close()
    var users []User
    json.NewDecoder(resp.Body).Decode(&users)
    c.JSON(200, users)
}

func AcceptDriver(c *gin.Context) {
    var input struct {
        Car   string `json:"car"`
        Plate string `json:"plate"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    userResp, _ := http.Get("http://localhost:8002/users/" + c.Param("id"))
    var user User
    json.NewDecoder(userResp.Body).Decode(&user)
    var roles []string
    json.Unmarshal([]byte(user.Roles), &roles)
    add := true
    for _, r := range roles { if r == "ROLE_DRIVER" { add = false; break } }
    if add { roles = append(roles, "ROLE_DRIVER") }
    rolesJSON, _ := json.Marshal(roles)
    patch := map[string]interface{}{
        "roles": string(rolesJSON),
        "car":   input.Car,
        "plate": input.Plate,
    }
    body, _ := json.Marshal(patch)
    req, _ := http.NewRequest("PATCH", "http://localhost:8002/users/"+c.Param("id"), bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    client := http.DefaultClient // Correction ici
    resp, err := client.Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": "PATCH failed"})
        return
    }
    defer resp.Body.Close()
    var updatedUser User
    json.NewDecoder(resp.Body).Decode(&updatedUser)
    c.JSON(200, updatedUser)
}

var driverRequests []DriverRequest

type DriverRequest struct {
    GoogleID  string `json:"user_id"`
    Car     string `json:"car"`
    Plate   string `json:"plate"`
    Status  string `json:"status"` 
}

func CreateDriverRequest(c *gin.Context) {
    var req DriverRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    req.Status = "pending"
    driverRequests = append(driverRequests, req)
    c.JSON(201, req)
}

func ListDriverRequests(c *gin.Context) {
    resp, err := http.Get("http://localhost:8002/driver-requests")
    if err != nil {
        c.JSON(500, gin.H{"error": "service-user injoignable"})
        return
    }
    defer resp.Body.Close()
    var reqs []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&reqs)
    c.JSON(200, reqs)
}

func HandleDriverRequest(c *gin.Context) {
    id := c.Param("id")
    var input struct{ Action string `json:"action"` }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(400, gin.H{"error": "bad input"})
        return
    }
    body, _ := json.Marshal(input)
    req, _ := http.NewRequest("PATCH", "http://localhost:8002/driver-requests/"+id, bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    client := http.DefaultClient // Correction ici
    resp, err := client.Do(req)
    if err != nil {
        c.JSON(500, gin.H{"error": "service-user KO"})
        return
    }
    defer resp.Body.Close()
    var updated map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&updated)
    c.JSON(200, updated)
}