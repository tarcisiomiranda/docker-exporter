from prometheus_client import generate_latest, Gauge, Info
from flask import Flask, Response
from datetime import datetime
from dateutil import parser
import argparse
import docker
import pytz

app = Flask(__name__)
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

@app.route('/metrics')
def metrics():
    get_container_metrics()
    return Response(generate_latest(), mimetype='text/plain')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--mode', default='prd', help='App mode: "dev" or "prd"', required=True)
    args = parser.parse_args()

    if args.mode.lower() == 'dev':
        print(' * Running: Exporter development...')
        app.run(debug=True, host='0.0.0.0', port=3003, use_reloader=True)

    elif args.mode.lower() == 'prd':
        print(' * Running: Exporter production...')
        app.run(debug=False, host='0.0.0.0', port=3003, use_reloader=False)

    else:
        print('Invalid app mode. Use "dev" or "prd".')
