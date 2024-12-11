package handlers

import (
	hash2 "PaymentSystem/iternal/jwt_auth/hash"
	"PaymentSystem/iternal/metrics"
	"PaymentSystem/iternal/users_db"
	"encoding/json"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func Register(db *gorm.DB) http.HandlerFunc {
	type RegUser struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestCount.Inc()
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(v)
		}))
		defer timer.ObserveDuration()
		w.Header().Set("Content-Type", "application/json")
		req := RegUser{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("failed to decode body")
			return
		}
		err = users_db.FindLoginFromDB(db, req.Login)
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("failed to find login from db")
			return
		}
		hash, err := hash2.HashPassword(req.Password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("failed to hash password")
			return
		}
		if err = users_db.AddToDB(db, req.Login, hash); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("failed to add to db")
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
