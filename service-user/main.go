package main

import (
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

var db *gorm.DB

func main() {
    godotenv.Load()
    dburl := os.Getenv("DATABASE_URL")
    if dburl == "" {
        panic("DATABASE_URL manquante !")
    }

    dsn := databaseURLToDSN(dburl)
    var err error
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Erreur connexion BDD: " + err.Error())
    }
    db.AutoMigrate(&User{})

    r := gin.Default()

r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
    r.GET("/users/:id", GetUserByID)
    r.GET("/users", ListUsers)         
    r.PATCH("/users/:id", PatchUser)   
    r.POST("/users/find-or-create", FindOrCreateUser)
    r.Run(":8002")
}

// Fonction qui convertit l'URL du .env vers un DSN Gorm/MySQL
func databaseURLToDSN(dburl string) string {
    dburl = strings.TrimPrefix(dburl, "mysql://")
    i := strings.Index(dburl, "@")
    creds := dburl[:i]
    after := dburl[i+1:]

    credsParts := strings.SplitN(creds, ":", 2)
    user, pass := credsParts[0], credsParts[1]

    i = strings.Index(after, "/")
    host := after[:i]
    rest := after[i+1:]

    i = strings.Index(rest, "?")
    dbname := rest
    params := ""
    if i != -1 {
        dbname = rest[:i]
        params = rest[i+1:]
    }
    dsn := user + ":" + pass + "@tcp(" + host + ")/" + dbname
    if params != "" {
        dsn += "?" + params
    }
    return dsn
}