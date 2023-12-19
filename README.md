# Docker Metrics Exporter for Prometheus

## Descrição
Este projeto fornece um exportador de métricas do Docker para o Prometheus, escrito em Python. Ele permite monitorar vários parâmetros de containers Docker, como tempo de atividade, status e imagem, expondo-os para o Prometheus.

## Características
- **Uptime do Container**: Mede o tempo desde o início de um container.
- **Status do Container**: Fornece o status atual do container.
- **Imagem do Container**: Informa a imagem usada pelo container.

## Requisitos
- Python 3
- Bibliotecas: `prometheus_client`, `flask`, `pytz`, `docker`
- Docker
- Prometheus

## Instalação e Uso
1. Clone o repositório: `git clone [URL do Repositório]`
2. Instale as dependências: `pip install -r requirements.txt`
3. Execute o script: `python docker_exporter.py --mode [dev|prd]`

## Endpoints
- `/metrics`: Retorna as métricas atuais dos containers.

## Modo de Execução
- Modo de Desenvolvimento: `python docker_exporter.py --mode dev`
- Modo de Produção: `python docker_exporter.py --mode prd`

## Contribuições
Contribuições são bem-vindas. Por favor, abra um issue ou pull request para discutir mudanças propostas.

## Licença
Este projeto está sob a Licença Pública Geral GNU (GPL), que é uma licença de software livre que garante a liberdade de compartilhar e alterar todo o software licenciado para garantir que ele permaneça livre.
