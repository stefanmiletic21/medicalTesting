package main

import (
	"math/rand"
	"net/http"
	"os"
	"medicalTesting/config"
	"medicalTesting/db"
	"medicalTesting/handler"
	"medicalTesting/logger"
	"time"

	"github.com/gorilla/context"
)

func main() {

	var err error

	rand.Seed(time.Now().UTC().UnixNano())

	logger.Init(os.Stdout, os.Stdout, os.Stderr)

	logger.Info("Initializing config")
	_, err = config.InitConfig("server", nil)
	if err != nil {
		logger.Error("Could not initialize config: %v\n", err)
		return
	}

	logger.Info("Initializing database")
	err = db.InitializeDb()
	if err != nil {
		logger.Error("Could not access database: %v\n", err)
		return
	}

	http.HandleFunc("/auth/", handler.HandleAuthorized)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/testingHash", handler.HandleTestingHash)

	logger.Info("Server listening on " + config.GetHTTPServerAddress())
	http.ListenAndServe(config.GetHTTPServerAddress(), context.ClearHandler(http.DefaultServeMux))

}
