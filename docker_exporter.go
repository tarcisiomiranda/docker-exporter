package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	containerUptime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_uptime_seconds",
			Help: "Time since container started in seconds.",
		},
		[]string{"container_name", "container_id"},
	)
	containerStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_status",
			Help: "Status of the container (1 for running, 0 for not running).",
		},
		[]string{"container_name", "container_id", "status"},
	)
	containerImage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_image",
			Help: "Image of the container.",
		},
		[]string{"container_name", "container_id", "image"},
	)
	port = flag.String("port", "3003", "Define a porta em que o servidor deve escutar")
)

func init() {
	prometheus.MustRegister(containerUptime)
	prometheus.MustRegister(containerStatus)
	prometheus.MustRegister(containerImage)
}

func recordMetrics() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		json, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			panic(err)
		}

		layout := "2006-01-02T15:04:05.999999999Z"
		startTime, err := time.Parse(layout, json.State.StartedAt)
		if err != nil {
			panic(err)
		}

		uptime := time.Since(startTime).Seconds()
		containerUptime.WithLabelValues(container.Names[0], container.ID).Set(uptime)

		status := 0
		if json.State.Running {
			status = 1
		}
		containerStatus.WithLabelValues(container.Names[0], container.ID, json.State.Status).Set(float64(status))

		containerImage.WithLabelValues(container.Names[0], container.ID, container.Image).Set(1)
	}
}

func main() {
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%s", *port)
	fmt.Printf("Starting server on http://%s/metrics\n", addr)

	go recordMetrics()

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
