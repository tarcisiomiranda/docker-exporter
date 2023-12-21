from prometheus_client import generate_latest, Gauge, Info, start_http_server
from datetime import datetime
from dateutil import parser
import argparse
import docker
import pytz
import time

client = docker.from_env()

container_uptime = Gauge('container_uptime_seconds', 'Time since container started', ['container_name', 'container_id'])
container_status = Gauge('container_status', 'Status of the container', ['container_name', 'container_id', 'status'])
container_image = Info('container_image', 'Image of the container', ['container_name', 'container_id'])

def get_container_metrics():
    containers = client.containers.list(all=True)
    for container in containers:
        container.reload()
        state = container.attrs['State']
        status = state['Status']
        image = container.image.tags[0] if container.image.tags else 'unknown'

        if 'StartedAt' in state and state['StartedAt'] != '0001-01-01T00:00:00Z':
            start_time = parser.parse(state['StartedAt'])
            uptime_seconds = (datetime.now(pytz.utc) - start_time).total_seconds()
            container_uptime.labels(container.name, container.short_id).set(uptime_seconds)

        container_status.labels(container.name, container.short_id, status).set(1)
        container_image.labels(container.name, container.short_id).info({'image': image})

if __name__ == '__main__':
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument('--port', default=3003, type=int, help='Port to listen on')
    args = arg_parser.parse_args()

    print('Starting Prometheus metrics server on port {}'.format(args.port))
    start_http_server(args.port)

    while True:
        get_container_metrics()
        time.sleep(1)
