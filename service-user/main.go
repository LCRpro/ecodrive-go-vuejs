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
    db.AutoMigrate(&DriverRequest{})
    db.AutoMigrate(&User{})
    db.AutoMigrate(&AppAccount{})

    var acc AppAccount
    if db.First(&acc).Error == gorm.ErrRecordNotFound {
        db.Create(&AppAccount{Balance: 0})
    }

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
    r.POST("/driver-requests", CreateDriverRequest)
    r.GET("/driver-requests", ListDriverRequests)
    r.PATCH("/driver-requests/:id", HandleDriverRequest)
    r.GET("/app-balance", GetAppBalance)
    r.Run(":8002")
}

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