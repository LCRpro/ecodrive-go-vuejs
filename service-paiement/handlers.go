package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)


var transactions []Transaction
var nextID uint
var txMu sync.Mutex

func saveTx(tx *Transaction) {
	txMu.Lock()
	defer txMu.Unlock()
	if tx.ID == 0 {
		tx.ID = nextID
		nextID++
	}
	transactions = append(transactions, *tx)
}

func findTxsByGoogleID(gid string) []Transaction {
	txMu.Lock()
	defer txMu.Unlock()
	var txs []Transaction
	for _, tx := range transactions {
		if tx.GoogleID == gid {
			txs = append(txs, tx)
		}
	}
	return txs
}

func allTxs() []Transaction {
	txMu.Lock()
	defer txMu.Unlock()
	out := make([]Transaction, len(transactions))
	copy(out, transactions)
	return out
}

var HasEnoughBalance = hasEnoughBalanceReal

func Deposit(c *gin.Context) {
	var input struct {
		GoogleID   string  `json:"google_id"`
		Amount     float64 `json:"amount"`
		CardNumber string  `json:"card_number"`
		Expiry     string  `json:"expiry"`
		CVV        string  `json:"cvv"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}
	if input.Amount < 10 {
		c.JSON(400, gin.H{"error": "Montant minimum 10€"})
		return
	}
	if input.GoogleID == "" {
		c.JSON(400, gin.H{"error": "Utilisateur inconnu"})
		return
	}
	tx := Transaction{
		GoogleID:  input.GoogleID,
		Type:      "deposit",
		Amount:    input.Amount,
		CreatedAt: time.Now(),
	}
	saveTx(&tx)

	go UpdateUserBalance(input.GoogleID, input.Amount)

	c.JSON(201, tx)
}

func Withdraw(c *gin.Context) {
	var input struct {
		GoogleID string  `json:"google_id"`
		Amount   float64 `json:"amount"`
		IBAN     string  `json:"iban"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "bad input"})
		return
	}

	if len(input.IBAN) < 8 {
		c.JSON(400, gin.H{"error": "IBAN incorrect"})
		return
	}

	ibanMask := input.IBAN[:2] + strings.Repeat("*", len(input.IBAN)-6) + input.IBAN[len(input.IBAN)-4:]

	if !HasEnoughBalance(input.GoogleID, input.Amount) {
		c.JSON(400, gin.H{"error": "Solde insuffisant"})
		return
	}

	tx := Transaction{
		GoogleID:  input.GoogleID,
		Type:      "withdrawal",
		Amount:    input.Amount,
		IBANMask:  ibanMask,
		CreatedAt: time.Now(),
	}
	saveTx(&tx)
	go UpdateUserBalance(input.GoogleID, -input.Amount)
	c.JSON(201, tx)
}

func ListTransactions(c *gin.Context) {
	gid := c.Param("google_id")
	txs := findTxsByGoogleID(gid)
	c.JSON(200, txs)
}

func ListAllTransactions(c *gin.Context) {
	c.JSON(200, allTxs())
}

func hasEnoughBalanceReal(googleID string, amount float64) bool {
	resp, err := http.Get(userServiceURL + "/users/" + googleID)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var user struct{ Balance float64 }
	json.NewDecoder(resp.Body).Decode(&user)
	return user.Balance >= amount
}

func UpdateUserBalance(googleID string, delta float64) {
	resp, err := http.Get(userServiceURL + "/users/" + googleID)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var user struct{ Balance float64 }
	json.NewDecoder(resp.Body).Decode(&user)
	newBalance := user.Balance + delta
	body, _ := json.Marshal(map[string]float64{"balance": newBalance})
	req, _ := http.NewRequest("PATCH", userServiceURL+"/users/"+googleID, strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Do(req)
}
