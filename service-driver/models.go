package main

import "time"

type Course struct {
	ID             uint64     `gorm:"primaryKey" json:"id"`
	PassengerID    string     `json:"passenger_id"`
	PassengerName  string     `json:"passenger_name"`
	PassengerEmail string     `json:"passenger_email"`
	StartLat       float64    `json:"start_lat"`
	StartLng       float64    `json:"start_lng"`
	EndLat         float64    `json:"end_lat"`
	EndLng         float64    `json:"end_lng"`
	DistanceKm     float64    `json:"distance_km"`
	CO2            float64    `json:"co2"`
	Amount         float64    `json:"amount"`
	Status         string     `json:"status"`
	DriverID       string     `json:"driver_id"`
	AcceptedAt     *time.Time `json:"accepted_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
