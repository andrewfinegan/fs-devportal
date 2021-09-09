// ################################################################################
// Copyright © 2021-2022 Fiserv, Inc. or its affiliates. 
// Fiserv is a trademark of Fiserv, Inc., 
// registered or used in the United States and foreign countries, 
// and may or may not be registered in your country.  
// All trademarks, service marks, 
// and trade names referenced in this 
// material are the property of their 
// respective owners. This work, including its contents 
// and programming, is confidential and its use 
// is strictly limited. This work is furnished only 
// for use by duly authorized licensees of Fiserv, Inc. 
// or its affiliates, and their designated agents 
// or employees responsible for installation or 
// operation of the products. Any other use, 
// duplication, or dissemination without the 
// prior written consent of Fiserv, Inc. 
// or its affiliates is strictly prohibited. 
// Except as specified by the agreement under 
// which the materials are furnished, Fiserv, Inc. 
// and its affiliates do not accept any liabilities 
// with respect to the information contained herein 
// and are not responsible for any direct, indirect, 
// special, consequential or exemplary damages 
// resulting from the use of this information. 
// No warranties, either express or implied, 
// are granted or extended by this work or 
// the delivery of this work
// ################################################################################


package main

import (
	"context"
	"devportal/api/product"
	"devportal/config"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
)

var cfg config.Config
var logger = logrus.New()

// Main method for the tenant server application
func main() {
	
	initApplication()
	//Configure server
	srv := &http.Server{
		Addr: "localhost:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      NewRouter(), // Pass our instance of gorilla/mux in.
	}

	logger.Info("Sample Tenant Server started on localhost port 8080")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Cant not start server due to error: ", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2000)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Info("Shutting down the server")
	os.Exit(0)

}

func initApplication() {
	initLogger()
	readConfig()
	setLoggingLevel()
}

func initLogger() {
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "2006/01/02 15:04:05.999999999Z07:00 MST"
	Formatter.FullTimestamp = true
	logger.SetFormatter(Formatter)

	config.Logger = logger
}

func setLoggingLevel() {
	switch config.AppConfig.Application.LoggingLevel {
	case "TraceLevel":
		logger.SetLevel(logrus.TraceLevel)
	case "DebugLevel":
		logger.SetLevel(logrus.DebugLevel)
	case "InfoLevel":
		logger.SetLevel(logrus.InfoLevel)
	case "WarnLevel":
		logger.SetLevel(logrus.WarnLevel)
	case "ErrorLevel":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	fmt.Printf("Application Logging Level: %v\n", logger.GetLevel())
}

func readConfig() {
	config.ReadFile(&cfg)
	config.ReadEnv(&cfg)

	cfg.GitHub.GitHubContentFullPath = cfg.GitHub.GitHubRawContentHost + "/" + cfg.GitHub.GitHubSourceOwner + "/" + cfg.GitHub.GitHubSourceRepo + "/" + cfg.GitHub.GitHubContentBranch + "/"
	product.DevPortalConfig = cfg

	logger.Debugf("%+v", cfg)
}
