package main

import (
	"PaymentSystem/config"
	"PaymentSystem/server"
	"PaymentSystem/storage"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	var config config.Config
	config.New()
	db, err := storage.NewDB()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to connect to database")
	}
	srv := server.Server{
		Rt:      mux.NewRouter(),
		Context: context.Background(),
		Db:      db,
		Config:  &config,
	}
	logrus.WithFields(logrus.Fields{}).Info("starting server")
	if err := srv.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		})
	}
}
