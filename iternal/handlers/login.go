package handlers

import (
	"PaymentSystem/iternal/jwt_auth/hash"
	"PaymentSystem/iternal/jwt_auth/jwt"
	"PaymentSystem/iternal/metrics"
	"PaymentSystem/iternal/users_db"
	"encoding/json"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	type User struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestCount.Inc()
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(v)
		}))
		defer timer.ObserveDuration()
		var user User
		w.Header().Add("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to decode body")
			return
		}
		hashPassword, err := users_db.GetHashPassword(db, user.Login)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get hash password")
			return
		}
		if err := hash.VerifyPassword(hashPassword, user.Password); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to verify password")
			return
		}
		token := jwt.GenerateJWT(user.Login)
		if token == "" {
			w.WriteHeader(http.StatusBadRequest)
			err = errors.New("token is empty")
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get token")
			return
		}
		_, err = w.Write([]byte(token))
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to write response")
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
