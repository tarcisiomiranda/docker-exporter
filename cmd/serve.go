package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	port           string
	updateInterval time.Duration
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Docker metrics exporter server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "9100", "Port to serve metrics on")
	serveCmd.Flags().DurationVarP(&updateInterval, "interval", "i", 15*time.Second, "Interval to update metrics (e.g., 15s, 1m)")
}

func startServer() {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Starting server on http://%s/metrics\n", addr)

	StartMetricsCollection(updateInterval)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
