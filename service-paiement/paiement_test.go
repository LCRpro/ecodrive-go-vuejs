package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
	"time"
)



func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&Transaction{})
	return db
}

var originalHasEnoughBalance func(string, float64) bool

func TestMain(m *testing.M) {
	originalHasEnoughBalance = HasEnoughBalance
	HasEnoughBalance = func(gid string, amt float64) bool { return true }
	m.Run()
	HasEnoughBalance = originalHasEnoughBalance
}

func TestDepositHandler(t *testing.T) {
	db = setupTestDB()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/deposit", Deposit)

	payload := map[string]interface{}{
		"google_id":   "user123",
		"amount":      50.0,
		"card_number": "1234567890123456",
		"expiry":      "12/29",
		"cvv":         "123",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("Expected 201, got %d (%s)", w.Code, w.Body.String())
	}
	var tx Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &tx); err != nil {
		t.Fatalf("Bad response: %v", err)
	}
	if tx.GoogleID != "user123" || tx.Amount != 50 || tx.Type != "deposit" {
		t.Errorf("Transaction fields not set as expected: %+v", tx)
	}
}

func TestDepositHandler_BadInput(t *testing.T) {
	db = setupTestDB()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/deposit", Deposit)

	payload := map[string]interface{}{
		"google_id": "",
		"amount":    5.0, // trop bas
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for bad input, got %d", w.Code)
	}
}

func TestWithdrawHandler(t *testing.T) {
	db = setupTestDB()
	gin.SetMode(gin.TestMode)
	HasEnoughBalance = func(gid string, amt float64) bool { return true }

	r := gin.Default()
	r.POST("/withdraw", Withdraw)

	payload := map[string]interface{}{
		"google_id": "user123",
		"amount":    15.5,
		"iban":      "FR7630006000011234567890189",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("Expected 201, got %d (%s)", w.Code, w.Body.String())
	}
	var tx Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &tx); err != nil {
		t.Fatalf("Bad response: %v", err)
	}
	if tx.Type != "withdrawal" {
		t.Errorf("Should have type withdrawal: %+v", tx)
	}
	if tx.IBANMask == "" || tx.IBANMask[:2] != "FR" {
		t.Errorf("Should mask IBAN: %+v", tx)
	}
}

func TestWithdrawHandler_BadIban(t *testing.T) {
	db = setupTestDB()
	gin.SetMode(gin.TestMode)
	HasEnoughBalance = func(gid string, amt float64) bool { return true }
	r := gin.Default()
	r.POST("/withdraw", Withdraw)
	payload := map[string]interface{}{
		"google_id": "user123",
		"amount":    20.0,
		"iban":      "1234", 
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for bad IBAN, got %d", w.Code)
	}
}

func TestListTransactions(t *testing.T) {
	db = setupTestDB()
	db.Create(&Transaction{GoogleID: "u1", Type: "deposit", Amount: 30, CreatedAt: time.Now()})
	db.Create(&Transaction{GoogleID: "u1", Type: "withdrawal", Amount: 10, CreatedAt: time.Now()})
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/transactions/:google_id", ListTransactions)
	req := httptest.NewRequest("GET", "/transactions/u1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var txs []Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &txs); err != nil {
		t.Fatalf("Could not decode transactions: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(txs))
	}
}

func TestListAllTransactions(t *testing.T) {
	db = setupTestDB()
	db.Create(&Transaction{GoogleID: "a", Type: "deposit", Amount: 80, CreatedAt: time.Now()})
	db.Create(&Transaction{GoogleID: "b", Type: "withdrawal", Amount: 21, CreatedAt: time.Now()})
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/transactions", ListAllTransactions)
	req := httptest.NewRequest("GET", "/transactions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var txs []Transaction
	if err := json.Unmarshal(w.Body.Bytes(), &txs); err != nil {
		t.Fatalf("Could not decode transactions: %v", err)
	}
	if len(txs) != 2 {
		t.Fatalf("Expected 2 transactions, got %d", len(txs))
	}
}