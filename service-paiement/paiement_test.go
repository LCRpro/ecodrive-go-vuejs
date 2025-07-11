package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
	"time"
)

func resetTestTxData() {
	transactions = nil
	nextID = 1
}

func TestDepositHandler(t *testing.T) {
	resetTestTxData()
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

func TestDepositHandler_BadJSON(t *testing.T) {
	resetTestTxData()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/deposit", Deposit)

	body := []byte(`{invalid}`)
	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestDepositHandler_BadAmount(t *testing.T) {
	resetTestTxData()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/deposit", Deposit)

	payload := map[string]interface{}{
		"google_id": "user123",
		"amount":    5.0,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for amount too low, got %d", w.Code)
	}
}

func TestDepositHandler_NoUser(t *testing.T) {
	resetTestTxData()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/deposit", Deposit)

	payload := map[string]interface{}{
		"amount": 15.0,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/deposit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for missing user, got %d", w.Code)
	}
}

func TestWithdrawHandler(t *testing.T) {
	resetTestTxData()
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

func TestWithdrawHandler_BadJSON(t *testing.T) {
	resetTestTxData()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/withdraw", Withdraw)
	body := []byte(`{invalid}`)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestWithdrawHandler_BadIban(t *testing.T) {
	resetTestTxData()
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

func TestWithdrawHandler_InsufficientBalance(t *testing.T) {
	resetTestTxData()
	gin.SetMode(gin.TestMode)
	HasEnoughBalance = func(gid string, amt float64) bool { return false }
	r := gin.Default()
	r.POST("/withdraw", Withdraw)
	payload := map[string]interface{}{
		"google_id": "user123",
		"amount":    20.0,
		"iban":      "FR7630006000011234567890189",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/withdraw", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("Expected 400 for insufficient balance, got %d", w.Code)
	}
}

func TestListTransactions(t *testing.T) {
	resetTestTxData()
	saveTx(&Transaction{GoogleID: "u1", Type: "deposit", Amount: 30, CreatedAt: time.Now()})
	saveTx(&Transaction{GoogleID: "u1", Type: "withdrawal", Amount: 10, CreatedAt: time.Now()})
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
	resetTestTxData()
	saveTx(&Transaction{GoogleID: "a", Type: "deposit", Amount: 80, CreatedAt: time.Now()})
	saveTx(&Transaction{GoogleID: "b", Type: "withdrawal", Amount: 21, CreatedAt: time.Now()})
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
