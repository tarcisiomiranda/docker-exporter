#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
  echo -e "${GREEN}[*]${NC} $1"
}

print_error() {
  echo -e "${RED}[!]${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}[!]${NC} $1"
}

if [ "$EUID" -ne 0 ]; then 
  print_error "Please run as root or with sudo"
  exit 1
fi

if ! command -v docker &> /dev/null; then
  print_warning "Docker is not installed. The exporter requires Docker to function."
  print_warning "Please install Docker before continuing."
  exit 1
fi

print_status "Starting Docker Exporter installation..."

if ! id "prometheus" &>/dev/null; then
  print_status "Creating prometheus user..."
  useradd -r -s /bin/false prometheus
fi

print_status "Adding prometheus user to docker group..."
usermod -aG docker prometheus

print_status "Creating directories..."
mkdir -p /opt/prometheus/docker_exporter
chown -R prometheus:prometheus /opt/prometheus/docker_exporter

print_status "Downloading Docker Exporter binary..."
curl -L https://github.com/tarcisiomiranda/docker-exporter/releases/download/v1.0.5/docker_exporter -o /opt/prometheus/docker_exporter/docker_exporter

if [ $? -ne 0 ]; then
  print_error "Failed to download binary"
  exit 1
fi

print_status "Setting permissions..."
chmod +x /opt/prometheus/docker_exporter/docker_exporter
chown prometheus:prometheus /opt/prometheus/docker_exporter/docker_exporter

print_status "Creating systemd service..."
cat << EOF > /etc/systemd/system/docker-exporter.service
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

print_status "Configuring systemd service..."
systemctl daemon-reload
systemctl enable docker-exporter
systemctl start docker-exporter

if systemctl is-active --quiet docker-exporter; then
  print_status "Docker Exporter is now running on port 9100"
  print_status "You can check the status with: systemctl status docker-exporter"
  print_status "View logs with: journalctl -u docker-exporter -f"
else
  print_error "Service failed to start. Please check logs with: journalctl -u docker-exporter -f"
  exit 1
fi

print_status "Installation completed!"
echo -e "\n${YELLOW}Remember to add this to your prometheus.yml:${NC}"
echo -e "
scrape_configs:
  - job_name: 'docker_exporter'
    static_configs:
      - targets: ['localhost:9100']
"
