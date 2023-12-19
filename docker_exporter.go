package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
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
	fmt.Println("Starting server on http://0.0.0.0:3003/metrics")
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe("0.0.0.0:3003", nil)
	if err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}
