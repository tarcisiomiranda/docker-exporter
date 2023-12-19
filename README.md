# Docker Exporter - Metrics for Prometheus

## Descrição
Este projeto fornece um exportador de métricas do Docker para o Prometheus, escrito em Python. Ele permite monitorar vários parâmetros de containers Docker, como tempo de atividade, status e imagem, expondo-os para o Prometheus.

## Características
- **Uptime do Container**: Mede o tempo desde o início de um container.
- **Status do Container**: Fornece o status atual do container.
- **Imagem do Container**: Informa a imagem usada pelo container.

## Requisitos Python
- Python 3
- Bibliotecas: `prometheus_client`, `flask`, `pytz`, `docker`
- Docker
- Prometheus

## Requisitos Go
Para compilar e executar a versão em Go do exportador:

```bash
go mod init docker_exporter
go get github.com/docker/docker/client
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

## Instalação e Uso (Python)
1. Clone o repositório: `git clone [URL do Repositório]`
2. Instale as dependências: `pip install -r requirements.txt`
3. Execute o script: `python docker_exporter.py --mode [dev|prd]`

## Instalação e Uso (Go) - Estático
1. Clone o repositório: `git clone [URL do Repositório]`
3. <s>Compile o código: GOOS=linux GOARCH=amd64 go build -o docker_exporter</s>
2. Compile o código: CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o docker_exporter .
3. Execute o binário: ./docker_exporter`


## Endpoints
- `/metrics`: Retorna as métricas atuais dos containers.

## Modo de Execução - Python
- Modo de Desenvolvimento: `python docker_exporter.py --mode dev`
- Modo de Produção: `python docker_exporter.py --mode prd`

## Modo de Execução - GO
- Modo de Desenvolvimento: `go run ./docker_exporter.go`
- Modo de Produção: `./docker_exporter`

## Instalação do binário GO no Linux

### Criando usuário *Prometheus* se setando as permissões
```
useradd prometheus
usermod -aG docker prometheus
mkdir /opt/prometheus/docker_exporter/
chown -R prometheus:prometheus /opt/prometheus/docker_exporter/
```

Conteúdo do services do systemd
***/etc/systemd/system/docker_exporter.service***
```
[Unit]
Description=Docker Exporter
Wants=network-online.target
After=network-online.target

# Internal Restart
StartLimitIntervalSec=600
StartLimitBurst=5

[Service]
User=prometheus
WorkingDirectory=/opt/prometheus/docker_exporter
ExecStart=/opt/prometheus/docker_exporter/docker_exporter
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target

```

Recarregue o systemd, habilite e de um start no exportador
```bash
systemctl daemon-reload
systemctl enable docker_exporter
systemctl start docker_exporter
```

Visualizar logs do programa em GO
```bash
journalctl -u docker_exporter
```

## Contribuições
Contribuições são bem-vindas. Por favor, abra um issue ou pull request para discutir mudanças propostas.

## Licença
Este projeto está sob a Licença Pública Geral GNU (GPL), que é uma licença de software livre que garante a liberdade de compartilhar e alterar todo o software licenciado para garantir que ele permaneça livre.
