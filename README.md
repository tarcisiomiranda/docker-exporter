# Docker Exporter - Metrics for Prometheus

## Description
This project provides a Docker metrics exporter for Prometheus, written in Go. It allows monitoring various Docker container parameters such as uptime, status, and image, exposing them to Prometheus.

> **Frontend Interface**: Check out our web interface at [docker-export-webui](https://github.com/tarcisiomiranda/docker-export-webui.git)

<!-- <img src="static/docker_exporter.png" width="800"/> -->
![Banner](static/docker_exporter.png)

## Quick Installation (Linux)
You can quickly install the latest version using this script:

```bash
#!/bin/bash

# Create directories and user
sudo useradd -r -s /bin/false prometheus
sudo usermod -aG docker prometheus
sudo mkdir -p /opt/prometheus/docker_exporter
sudo chown -R prometheus:prometheus /opt/prometheus/docker_exporter

# Download the binary
sudo curl -L https://github.com/tarcisiomiranda/docker-exporter/releases/download/v1.0.5/docker_exporter -o /opt/prometheus/docker_exporter/docker_exporter

# Make it executable
sudo chmod +x /opt/prometheus/docker_exporter/docker_exporter

# Create systemd service file
cat << EOF | sudo tee /etc/systemd/system/docker-exporter.service
[Unit]
Description=Docker Exporter
Wants=network-online.target
After=network-online.target docker.service

StartLimitIntervalSec=600
StartLimitBurst=5

[Service]
User=prometheus
WorkingDirectory=/opt/prometheus/docker_exporter
ExecStart=/opt/prometheus/docker_exporter/docker_exporter serve -p 9100
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload
sudo systemctl enable docker-exporter
sudo systemctl start docker-exporter

echo "Installation completed! Docker Exporter is running on port 9100"
```

Save this script as `install-docker-exporter.sh` and run:
```bash
curl -O https://raw.githubusercontent.com/tarcisiomiranda/docker-exporter/main/install-docker-exporter.sh
chmod +x install-docker-exporter.sh
sudo ./install-docker-exporter.sh
```

## Features
- **Container Uptime**: Measures the time since a container started.
- **Container Status**: Provides the current container status.
- **Container Image**: Shows the image used by the container.

## Development Setup
To compile and run the Docker Exporter from source:

```bash
# Clone the repository
git clone https://github.com/tarcisiomiranda/docker-exporter.git
cd docker-exporter

# Initialize Go modules
go mod tidy

# Build the binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o docker_exporter .
```

## Endpoints
- `/metrics`: Returns current container metrics (default port: 9100)

## Execution Modes
- Development Mode: `go run docker_exporter.go serve -p 9100`
- Production Mode: `./docker_exporter serve -p 9100`

## Monitoring
View service logs:
```bash
journalctl -u docker-exporter -f
```

Check service status:
```bash
systemctl status docker-exporter
```

## Building from Source
Requirements:
- Go 1.22 or later
- Docker (for monitoring containers)

## Prometheus Configuration
Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'docker_exporter'
    static_configs:
      - targets: ['localhost:9100']
```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the GNU General Public License (GPL).

## Support
If you encounter any issues or have questions, please open an issue on GitHub.
