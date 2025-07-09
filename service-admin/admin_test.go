package main

import (
     "os"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "time"
    "strings"
    "errors"
)

func init() {
    os.Setenv("JWT_SECRET", "devsecret")
}
type mockRoundTripper struct {
    roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.roundTripFunc(req)
}

func setMockClient(fn func(*http.Request) (*http.Response, error)) func() {
    origClient := http.DefaultClient
    http.DefaultClient = &http.Client{
        Transport: &mockRoundTripper{roundTripFunc: fn},
    }
    return func() { http.DefaultClient = origClient }
}

// -------- JWT UTILS --------

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

// --------- TESTS ---------

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

func TestListUsers_Success(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        users := []User{{ID: 1, Name: "Toto"}}
        buf, _ := json.Marshal(users)
        return &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(bytes.NewReader(buf)),
            Header:     make(http.Header),
        }, nil
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    ListUsers(c)
    if w.Code != 200 {
        t.Fatalf("Expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "Toto") {
        t.Error("User not in response")
    }
}

func TestListUsers_Error(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        return nil, errors.New("fail")
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    ListUsers(c)
    if w.Code != 500 {
        t.Fatalf("Expected 500, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "injoignable") {
        t.Error("Expected error message for injoignable service")
    }
}

func TestAcceptDriver_Success(t *testing.T) {
    // Mock user GET and PATCH
    callCount := 0
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        callCount++
        if req.Method == "GET" {
            // Renvoie un user sans ROLE_DRIVER
            u := User{ID: 1, Name: "Liam", Roles: `["ROLE_PASSAGER"]`}
            buf, _ := json.Marshal(u)
            return &http.Response{
                StatusCode: 200,
                Body:       ioutil.NopCloser(bytes.NewReader(buf)),
                Header:     make(http.Header),
            }, nil
        }
        if req.Method == "PATCH" {
            // Retourne le user mis à jour
            u := User{ID: 1, Name: "Liam", Roles: `["ROLE_PASSAGER","ROLE_DRIVER"]`, Car: "Tesla", Plate: "AB-123-CD"}
            buf, _ := json.Marshal(u)
            return &http.Response{
                StatusCode: 200,
                Body:       ioutil.NopCloser(bytes.NewReader(buf)),
                Header:     make(http.Header),
            }, nil
        }
        return nil, errors.New("unexpected")
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/accept/1", bytes.NewBuffer([]byte(`{"car":"Tesla","plate":"AB-123-CD"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    AcceptDriver(c)
    if w.Code != 200 {
        t.Fatalf("Expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "ROLE_DRIVER") {
        t.Error("ROLE_DRIVER not in response")
    }
}

func TestAcceptDriver_BadInput(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/accept/1", bytes.NewBuffer([]byte(`{"car":""`))) // Malformed JSON
    c.Request.Header.Set("Content-Type", "application/json")
    AcceptDriver(c)
    if w.Code != 400 {
        t.Fatalf("Expected 400 for bad input, got %d", w.Code)
    }
}

func TestAcceptDriver_PatchError(t *testing.T) {
    // GET user OK, PATCH KO
    callCount := 0
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        callCount++
        if req.Method == "GET" {
            u := User{ID: 1, Name: "Liam", Roles: `["ROLE_PASSAGER"]`}
            buf, _ := json.Marshal(u)
            return &http.Response{
                StatusCode: 200,
                Body:       ioutil.NopCloser(bytes.NewReader(buf)),
                Header:     make(http.Header),
            }, nil
        }
        if req.Method == "PATCH" {
            return nil, errors.New("fail patch")
        }
        return nil, errors.New("unexpected")
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/accept/1", bytes.NewBuffer([]byte(`{"car":"Tesla","plate":"AB-123-CD"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    AcceptDriver(c)
    if w.Code != 500 {
        t.Fatalf("Expected 500 on PATCH fail, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "PATCH failed") {
        t.Error("Expected PATCH failed error")
    }
}

func TestCreateDriverRequest_Success(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/driver-request", bytes.NewBuffer([]byte(`{"user_id":"abc123","car":"Toyota","plate":"ZZ-999-YY"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    CreateDriverRequest(c)
    if w.Code != 201 {
        t.Fatalf("Expected 201, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), `"status":"pending"`) {
        t.Error("Driver request should be pending")
    }
}

func TestCreateDriverRequest_BadInput(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/driver-request", bytes.NewBuffer([]byte(`notjson`)))
    c.Request.Header.Set("Content-Type", "application/json")
    CreateDriverRequest(c)
    if w.Code != 400 {
        t.Fatalf("Expected 400 for bad input, got %d", w.Code)
    }
}

func TestListDriverRequests_Success(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        reqs := []map[string]interface{}{{"user_id": "abc", "status": "pending"}}
        buf, _ := json.Marshal(reqs)
        return &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(bytes.NewReader(buf)),
            Header:     make(http.Header),
        }, nil
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    ListDriverRequests(c)
    if w.Code != 200 {
        t.Fatalf("Expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "pending") {
        t.Error("Driver request status pending missing")
    }
}

func TestListDriverRequests_Error(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        return nil, errors.New("fail")
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    ListDriverRequests(c)
    if w.Code != 500 {
        t.Fatalf("Expected 500, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "injoignable") {
        t.Error("Expected injoignable service error")
    }
}

func TestHandleDriverRequest_Success(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        updated := map[string]interface{}{"id": "1", "status": "approved"}
        buf, _ := json.Marshal(updated)
        return &http.Response{
            StatusCode: 200,
            Body:       ioutil.NopCloser(bytes.NewReader(buf)),
            Header:     make(http.Header),
        }, nil
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/driver-requests/1", bytes.NewBuffer([]byte(`{"action":"approve"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    HandleDriverRequest(c)
    if w.Code != 200 {
        t.Fatalf("Expected 200, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "approved") {
        t.Error("Expected approved status in response")
    }
}

func TestHandleDriverRequest_BadInput(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/driver-requests/1", bytes.NewBuffer([]byte(`notjson`)))
    c.Request.Header.Set("Content-Type", "application/json")
    HandleDriverRequest(c)
    if w.Code != 400 {
        t.Fatalf("Expected 400 for bad input, got %d", w.Code)
    }
}

func TestHandleDriverRequest_ServiceError(t *testing.T) {
    restore := setMockClient(func(req *http.Request) (*http.Response, error) {
        return nil, errors.New("fail")
    })
    defer restore()

    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = gin.Params{{Key: "id", Value: "1"}}
    c.Request = httptest.NewRequest("PATCH", "/driver-requests/1", bytes.NewBuffer([]byte(`{"action":"reject"}`)))
    c.Request.Header.Set("Content-Type", "application/json")
    HandleDriverRequest(c)
    if w.Code != 500 {
        t.Fatalf("Expected 500 for service KO, got %d", w.Code)
    }
    if !strings.Contains(w.Body.String(), "service-user KO") {
        t.Error("Expected service-user KO error")
    }
}