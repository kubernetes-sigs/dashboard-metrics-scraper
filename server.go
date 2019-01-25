package main

import (
	"database/sql"
	"flag"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	sideapi "github.com/jeefy/metrics-sidecar/pkg/api"
	sidedb "github.com/jeefy/metrics-sidecar/pkg/database"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var kubeconfig *string
	var dbFile *string
	var refreshInterval *int
	var maxWindow *int

	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	dbFile = flag.String("db-file", ":memory:", "What file to use as a SQLite3 database. Defaults to ':memory:'")
	refreshInterval = flag.Int("refresh-interval", 10, "Frequency (in seconds) to update the metrics database. Defaults to '5'")
	maxWindow = flag.Int("max-window", 15, "Window of time you wish to retain records (in minutes). Defaults to '15'")

	flag.Parse()

	// This should only be run in-cluster so...
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Unable to generate a client config: %s", err)
	}

	log.Infof("Kubernetes host: %s", config.Host)

	// Generate the metrics client
	clientset, err := metricsclient.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to generate a clientset: %s", err)
	}

	// Create the db "connection"
	db, err := sql.Open("sqlite3", *dbFile)
	if err != nil {
		log.Fatalf("Unable to open Sqlite database: %s", err)
	}
	defer db.Close()

	// Populate tables
	err = sidedb.CreateDatabase(db)
	if err != nil {
		log.Fatalf("Unable to initialize database tables: %s", err)
	}

	go func() {
		r := mux.NewRouter()
		sideapi.ApiManager(r, db)
		// Bind to a port and pass our router in
		log.Fatal(http.ListenAndServe(":8000", r))
	}()

	// Start the machine. Scrape every refreshInterval
	ticker := time.NewTicker(time.Duration(*refreshInterval) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-quit:
			ticker.Stop()
			return

		case <-ticker.C:
			err = nil
			nodeMetrics, err := clientset.Metrics().NodeMetricses().List(v1.ListOptions{})
			if err != nil {
				log.Errorf("Error scraping node metrics: %s", err)
				break
			}

			podMetrics, err := clientset.Metrics().PodMetricses("").List(v1.ListOptions{})
			if err != nil {
				log.Errorf("Error scraping pod metrics: %s", err)
				break
			}

			// Insert scrapes into DB
			err = sidedb.UpdateDatabase(db, nodeMetrics, podMetrics)
			if err != nil {
				log.Errorf("Error updating database: %s", err)
				break
			}

			// Delete rows outside of the maxWindow time
			err = sidedb.CullDatabase(db, maxWindow)
			if err != nil {
				log.Errorf("Error culling database: %s", err)
				break
			}

			log.Info("Database updated")
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
