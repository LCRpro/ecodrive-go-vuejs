package main

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db.AutoMigrate(&Transaction{})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ecodrive.liamcariou.fr"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/deposit", Deposit)
	r.POST("/withdraw", Withdraw)
	r.GET("/transactions/:google_id", ListTransactions)
	r.GET("/transactions", ListAllTransactions)
	r.Run(":8004")
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
