package main

import (
	"database/sql"
	"example/web-service-gin/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"time"
)

func init() {
	utils.PopulateConfig("credentials.json")
}

func main() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		utils.Cred.Host, utils.Cred.Port, utils.Cred.User, utils.Cred.Password, utils.Cred.Dbname)
	utils.Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(utils.Db)

	err = utils.Db.Ping()
	if err != nil {
		panic(err)
	}

	server := gin.Default()
	server.SetTrustedProxies([]string{"192.168.1.7", "192.168.1.4"})

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://localhost:4200"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			"accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router := server.Group("/api")
	addValidators()
	addAPIRoutes(router)

	go utils.EmptyVerificationCodesRoutine(10 * 60)
	go utils.ClearExpiredSessions(24 * 60 * 60)

	utils.MutexSession.Lock()
	err = utils.GetSessionsAfterRestart(utils.SessionToEmailID)
	utils.MutexSession.Unlock()

	if err != nil {
		panic(err)
	}

	err = server.RunTLS("0.0.0.0:8080", "localhost.crt", "localhost.key")
	if err != nil {
		return
	}
}
