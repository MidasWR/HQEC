package handlers

import (
	"PaymentSystem/iternal/jwt_auth/jwt"
	kafka2 "PaymentSystem/iternal/kafka"
	"PaymentSystem/iternal/metrics"
	"PaymentSystem/iternal/users_db"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func AnalyticsJSONPage(db *gorm.DB) http.HandlerFunc {
	type Response struct {
		Date    time.Time `json:"date"`
		Balance float32   `json:"balance"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestCount.Inc()
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(v)
		}))
		defer timer.ObserveDuration()
		w.Header().Add("Content-Type", "application/json")
		key := r.Header.Get("apikey")
		logrus.Printf("Analytics key: %s", key)
		login, err := jwt.ParseJWT(key)
		logrus.WithFields(logrus.Fields{
			"login": login,
		}).Info("Parsed JWT")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to parse JWT")
			return
		}
		id, err := users_db.GetIdByLogin(db, login)
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Info("Get Id")
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get id from db")
			return
		}
		ids := strconv.Itoa(id)
		consumer := kafka2.CreateConsumer()
		messages, err := kafka2.ConsumeStructuredMessages(consumer, "analytics"+ids)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to consume messages")
		}
		var responses []Response
		for _, message := range messages {
			var response Response
			response.Date = message.Date
			response.Balance = message.Balance
			responses = append(responses, response)
		}
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to encode responses to json")
		} else {
			logrus.WithFields(logrus.Fields{
				"id": id,
			}).Info("Successfully processed messages")
		}

		w.WriteHeader(http.StatusOK)

	}
}
