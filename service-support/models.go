package main

import "time"

type SupportTicket struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    GoogleID   string    `json:"google_id"`
    Category   string    `json:"category"`
    Message    string    `json:"message"`
    Status     string    `json:"status"`
    AdminReply string    `json:"admin_reply,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}