package main

import (
	"database/sql"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"

	sideapi "github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/api"
	sidedb "github.com/kubernetes-sigs/dashboard-metrics-scraper/pkg/database"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var kubeconfig *string
	var dbFile *string
	var metricResolution *time.Duration
	var metricDuration *time.Duration

	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	kubeconfig = flag.String("kubeconfig", "", "The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	dbFile = flag.String("db-file", ":memory:", "What file to use as a SQLite3 database.")
	metricResolution = flag.Duration("metric-resolution", 60 * time.Second, "The resolution at which dashboard-metrics-scraper will poll metrics.")
	metricDuration = flag.Duration("metric-duration", 15 * time.Minute, "The duration after which metrics are purged from the database.")

	flag.Set("logtostderr", "true")
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

	// Start the machine. Scrape every metricResolution
	ticker := time.NewTicker(*metricResolution)
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
			err = sidedb.CullDatabase(db, metricDuration)
			if err != nil {
				log.Errorf("Error culling database: %s", err)
				break
			}

			log.Infof("Database updated: %d nodes, %d pods", len(nodeMetrics.Items), len(podMetrics.Items))
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
