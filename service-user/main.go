package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

var db *gorm.DB

func main() {
    dsn := "root:root@tcp(127.0.0.1:3306)/goecodrive?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    db.AutoMigrate(&User{})
}