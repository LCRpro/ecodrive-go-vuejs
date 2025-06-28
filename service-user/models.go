package main

import "gorm.io/gorm"

type User struct {
    gorm.Model
    GoogleID  string `gorm:"uniqueIndex"`
    Email     string `gorm:"uniqueIndex"`
    Name      string
    Roles     string
    Birthdate string
    Gender    string
    Car       string
    Plate     string
}