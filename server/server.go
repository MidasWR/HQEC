package server

import (
	"PaymentSystem/config"
	_ "PaymentSystem/docs"
	handlers2 "PaymentSystem/iternal/handlers"
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
	"net/http"
)

// @title           Payment System API
// @version         1.0
// @description     This is a sample server for Payment System
// @host            localhost:8080
// @schemes         http https

type Server struct {
	Db      *gorm.DB
	Context context.Context
	Rt      *mux.Router
	Config  *config.Config
}

func (s *Server) Run() error {
	RegistrationHandle(s)
	LoginHandle(s)
	TransactionsHandleGet(s)
	TransactionsHandlePost(s)
	BalanceHandle(s)
	MetricsHandle(s)
	s.Rt.Handle("/swagger/", httpSwagger.WrapHandler)
	return http.ListenAndServe(s.Config.Host+":"+s.Config.Port, s.Rt)
}

// @Summary Registration Page
// @Description Page where you can register in the server's database
// @Tags Authentication/Authorization
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /registration [post]
func RegistrationHandle(s *Server) {
	s.Rt.HandleFunc("/registration", handlers2.Register(s.Db)).Methods("POST")
}

// @Summary Login Page
// @Description Page where you can log into the system
// @Tags Authentication/Authorization
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "bearer_token"
// @Router /login [post]
func LoginHandle(s *Server) {
	s.Rt.HandleFunc("/login", handlers2.LoginHandler(s.Db)).Methods("POST")
}

// @Summary List of Transactions
// @Description Page where you can get a list of transactions for the user
// @Tags Logic
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /transactions [get]
func TransactionsHandleGet(s *Server) {
	s.Rt.HandleFunc("/transactions", handlers2.GetTransactions(s.Db)).Methods("GET")
}

// @Summary Add Transaction
// @Description Page where you can add a new transaction to the list
// @Tags Logic
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /transactions [post]
func TransactionsHandlePost(s *Server) {
	s.Rt.HandleFunc("/transactions", handlers2.PostTransactions(s.Db)).Methods("POST")
}

// @Summary Get Balance
// @Description Page where you can get the balance for a client
// @Tags Logic
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /balance [get]
func BalanceHandle(s *Server) {
	s.Rt.HandleFunc("/balance", handlers2.GetBalance(s.Db)).Methods("GET")
}

// @Summary Metrics Page
// @Description Page for fetching custom metrics about the server's operation
// @Tags Metrics
// @Accept application/json
// @Produce application/json
// @Success 200
// @Router /metrics [get]
func MetricsHandle(s *Server) {
	s.Rt.Handle("/metrics", promhttp.Handler())
}
