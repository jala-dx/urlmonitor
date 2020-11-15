## URL Monitor
Simple URL Monitor to check the connectivity of external endpoints.
This service listens on port 2112 at /metrics prefix

### How to deploy
- Refer to the docker image at do https://hub.docker.com/repository/docker/jalaja/urlimg
  - docker pull jalaja/urlimg:1.0
- Copy the Podspec (an easy environment can be kind)
   - kind create cluster
   - kubectl cluster-info --context kind-kin
   - kubectly apply -f urlmonitor_k8.yaml
- Docker pull prometheus images



#### Console Logs

```bash
>kubectl get all -n urlns
NAME                              READY   STATUS    RESTARTS   AGE
pod/urlmonitor-76d599cdbf-f49tq   1/1     Running   0          3m56s

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/urlmonitor   ClusterIP   10.97.102.157   <none>        2112/TCP   6m23s

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/urlmonitor   1/1     1            1           6m23s

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/urlmonitor-76d599cdbf   1         1         1       3m56s


root@urlmonitor-76d599cdbf-f49tq:/# curl http://10.97.102.157:2112/metrics
# HELP external_url_up Connectivity status of an endpoint url.
# TYPE external_url_up counter
external_url_up{url="https://httpstat.us/503"} 0
# HELP external_url_response_ms Latency to reach endpoint url.
# TYPE external_url_response_ms counter
external_url_response_ms{url="https://httpstat.us/503"} 343
# HELP external_url_up Connectivity status of an endpoint url.
# TYPE external_url_up counter
external_url_up{url="https://httpstat.us/200"} 1
# HELP external_url_response_ms Latency to reach endpoint url.
# TYPE external_url_response_ms counter
external_url_response_ms{url="https://httpstat.us/200"} 19
# HELP external_url_up Connectivity status of an endpoint url.
# TYPE external_url_up counter
external_url_up{url="http://golangcode.com/robots.txt"} 1
# HELP external_url_response_ms Latency to reach endpoint url.
# TYPE external_url_response_ms counter
external_url_response_ms{url="http://golangcode.com/robots.txt"} 235
root@urlmonitor-76d599cdbf-f49tq:/# 

```

#### Prometheus connection to the urlmonitor service


```bash
- docker pull prom/prometheus:latest

- docker images | grep prom
prom/prometheus                   latest              7cc97b58fb0e        9 days ago          168.3 MB

- docker run  --name prometheus -p 9090:9090 -v /local/jganapat/prom/config/prom.yml:/etc/prometheus/prometheus.yml prom/prometheus --config.file=/etc/prometheus/prometheus.yml
```

#### Prometheus configuration 

```bash
prom.yml

# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.
    static_configs:
    - targets: ['127.0.0.1:9090']

  - job_name: 'myapp'
    #metrics_path: '/metrics'
    scrape_interval: 5s
    static_configs:
    - targets: ['x.x.x.x:2112'] == Make sure to specify the correct IP address here to talk to the urlmonitor service

```
