package main

type User struct {
    ID        uint   `gorm:"primaryKey" json:"id"`
    GoogleID  string `gorm:"size:191;uniqueIndex" json:"google_id"`
    Email     string `gorm:"size:191;uniqueIndex" json:"email"`
    Name      string `gorm:"size:191" json:"name"`
    Roles     string `gorm:"type:text" json:"roles"`
    Birthdate string `gorm:"size:191" json:"birthdate"`
    Gender    string `gorm:"size:191" json:"gender"`
    Car       string `gorm:"size:191" json:"car"`
    Plate     string `gorm:"size:191" json:"plate"`
}