package handlers

import (
	"PaymentSystem/iternal/balance_db"
	"PaymentSystem/iternal/brokers/kafka"
	"PaymentSystem/iternal/jwt_auth/jwt"
	"PaymentSystem/iternal/metrics"
	"PaymentSystem/iternal/users_db"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetBalance(db *gorm.DB) http.HandlerFunc {
	type Balance struct {
		Id      int       `gorm:"primaryKey"`
		UserId  int       `gorm:"foreignKey:UserId;references:Id"`
		Balance float32   `gorm:"type:decimal(10,2)"`
		Date    time.Time `gorm:"type:date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestCount.Inc()
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(v)
		}))
		defer timer.ObserveDuration()
		w.Header().Add("Content-Type", "application/json")
		key := r.URL.Query().Get("apikey")
		login, err := jwt.ParseJWT(key)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to parse JWT")
			return
		}
		id, err := users_db.GetIdByLogin(db, login)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get id from db")
			return
		}
		balance, err := balance_db.GetBalance(db, id)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get balance from db")
			return
		}
		if _, err := fmt.Fprint(w, balance); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to write response")
			return
		}
		bls, err := balance_db.GetBalanceStructById(db, id)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Debugln("Failed to get balance struct from db")
		}
		producer := kafka.CreateProducer()
		defer producer.Close()
		topic := "balance"
		msg := kafka.KafkaMessage{
			Id:      bls.Id,
			Balance: bls.Balance,
			UserId:  bls.UserId,
			Date:    bls.Date,
		}
		kafka.ProduceMessage(producer, topic, msg)
		w.WriteHeader(http.StatusOK)
	}
}
