package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func init() {
	PopulateConfig("credentials.json")
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cred.Host, cred.Port, cred.User, cred.Password, cred.Dbname)
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

	CORSMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		}
	}

	server.Use(CORSMiddleware())

	router := server.Group("/api")
	addValidators()
	addAPIRoutes(router)

	go emptyVerificationCodesRoutine(10 * 60)
	go clearExpiredSessions(24 * 60 * 60)

	err = server.Run("localhost:8080")
	if err != nil {
		return
	}
}
