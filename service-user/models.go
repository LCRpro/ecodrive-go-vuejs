package main

type User struct {
    ID        uint    `gorm:"primaryKey" json:"id"`
    GoogleID  string  `gorm:"size:191;uniqueIndex" json:"google_id"`
    Email     string  `gorm:"size:191;uniqueIndex" json:"email"`
    Name      string  `gorm:"size:191" json:"name"`
    Roles     string  `gorm:"type:text" json:"roles"`
    Birthdate string  `gorm:"size:191" json:"birthdate"`
    Gender    string  `gorm:"size:191" json:"gender"`
    Car       string  `gorm:"size:191" json:"car"`
    Plate     string  `gorm:"size:191" json:"plate"`
    Balance   float64 `gorm:"default:0" json:"balance"` 
}

type AppAccount struct {
    ID      uint    `gorm:"primaryKey"`
    Balance float64 `json:"balance"` 
}

func (AppAccount) TableName() string { return "app_account" }

type DriverRequest struct {
    ID      uint   `gorm:"primaryKey"`
    GoogleID string `json:"google_id"`
    Car     string `json:"car"`
    Plate   string `json:"plate"`
    Status  string `json:"status"` 
}

func (DriverRequest) TableName() string { return "driver_requests" }