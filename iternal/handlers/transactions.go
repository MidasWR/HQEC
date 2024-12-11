package handlers

import (
	"PaymentSystem/iternal/balance_db"
	"PaymentSystem/iternal/jwt_auth/jwt"
	"PaymentSystem/iternal/metrics"
	"PaymentSystem/iternal/transactions_db"
	"PaymentSystem/iternal/users_db"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func GetTransactions(db *gorm.DB) http.HandlerFunc {
	type Transaction struct {
		Id    int       `gorm:"primaryKey"`
		Date  time.Time `gorm:"type:datetime"`
		Money float32   `gorm:"type:decimal(10,2)"`
	}
	type Params struct {
		DateStart *time.Time `json:"date_start"`
		DateEnd   *time.Time `json:"date_end"`
		Date      *time.Time `json:"date"`
		Money     *float32   `json:"money"`
		Page      int        `json:"page"`
		PerPage   int        `json:"per_page"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.RequestCount.Inc()
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
		}
		params := Params{}
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to decode request body")
			return
		}
		transactions, err := transactions_db.GetTransactionsFromDb(transactions_db.Params(params), db, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get transactions from db")
			return
		}
		for _, transaction := range transactions {
			fmt.Println(transaction)
		}
		w.WriteHeader(http.StatusOK)
	}
}
func PostTransactions(db *gorm.DB) http.HandlerFunc {
	type Transactions struct {
		Id     int       `gorm:"primaryKey"`
		UserId int       `gorm:"index"`
		Date   time.Time `gorm:"type:datetime"`
		Money  float32   `gorm:"type:decimal(10,2)"`
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
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to get id from db")
			return
		}
		trn := Transactions{}
		err = json.NewDecoder(r.Body).Decode(&trn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to decode request body")
			return
		}
		if err := transactions_db.AddTransactionToDb(transactions_db.Transactions(trn), db, id); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to add transaction to db")
			return
		}

		if err := balance_db.ReBalance(db, trn.Money, id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Warn("Failed to re balance")
			return
		}
		
		w.WriteHeader(http.StatusCreated)
	}
}
