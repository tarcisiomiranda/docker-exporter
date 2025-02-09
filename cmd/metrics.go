package cmd

import (
	"context"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
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

func StartMetricsCollection(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			updateMetrics()
			<-ticker.C
		}
	}()
}

func updateMetrics() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Printf("Error creating Docker client: %v", err)
		return
	}
	defer cli.Close()

	containerUptime.Reset()
	containerStatus.Reset()
	containerImage.Reset()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return
	}

	for _, container := range containers {
		json, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Printf("Error inspecting container %s: %v", container.ID, err)
			continue
		}

		containerName := container.Names[0]
		if len(containerName) > 0 && containerName[0] == '/' {
			containerName = containerName[1:]
		}

		layout := "2006-01-02T15:04:05.999999999Z"
		startTime, err := time.Parse(layout, json.State.StartedAt)
		if err != nil {
			log.Printf("Error parsing start time for container %s: %v", container.ID, err)
			continue
		}

		uptime := time.Since(startTime).Seconds()
		containerUptime.WithLabelValues(containerName, container.ID).Set(uptime)

		status := 0
		if json.State.Running {
			status = 1
		}
		containerStatus.WithLabelValues(containerName, container.ID, json.State.Status).Set(float64(status))

		containerImage.WithLabelValues(containerName, container.ID, container.Image).Set(1)
	}
}
