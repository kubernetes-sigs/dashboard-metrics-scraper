package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	dashboardProvider "github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/api/dashboard"
	_ "github.com/mattn/go-sqlite3"
)

func ApiManager(r *mux.Router, db *sql.DB) {
	dashboardRouter := r.PathPrefix("/api/v1/dashboard").Subrouter()
	dashboardProvider.DashboardRouter(dashboardRouter, db)
	r.PathPrefix("/").HandlerFunc(DefaultHandler)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("URL: %s", r.URL)
	w.Write([]byte(msg))
}
