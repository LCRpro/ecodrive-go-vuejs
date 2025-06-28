package main

import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestUserCreate(t *testing.T) {
    testdb, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    testdb.AutoMigrate(&User{})
    user := User{GoogleID: "id", Email: "test@example.com", Name: "Test"}
    if err := testdb.Create(&user).Error; err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }
}