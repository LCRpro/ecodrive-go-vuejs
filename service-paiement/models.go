package main

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GoogleID  string    `json:"google_id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	IBANMask  string    `json:"iban_mask,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
