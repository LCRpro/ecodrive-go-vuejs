package main

type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    GoogleID  string `gorm:"uniqueIndex" json:"google_id"`
    Email     string `gorm:"uniqueIndex" json:"email"`
    Name      string `json:"name"`
    Roles     string `json:"roles"`
    Birthdate string `json:"birthdate"`
    Gender    string `json:"gender"`
    Car       string `json:"car"`
    Plate     string `json:"plate"`
}