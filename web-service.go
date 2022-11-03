package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

const (
	host     = "localhost"
	port     = 8184
	user     = "postgres"
	password = "root"
	dbname   = "godb"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.SetTrustedProxies([]string{"192.168.1.2"})

	router := server.Group("/api")
	addValidators()
	addAPIRoutes(router)

	err = server.Run("localhost:8080")
	if err != nil {
		return
	}
}
