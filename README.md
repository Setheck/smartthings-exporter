# smartthings-exporter

A simple app to help monitor your smart devices via prometheus. 

## Run

image is available for both linux/amd64 and linux/arm/v7 [here](https://hub.docker.com/r/setheck/smartthings-exporter)

run with docker
```
docker run --rm \
    -p "9119:9119" \
    -e "STE_API_TOKEN=<your api token>" \
    setheck/smartthings-exporter
```

## Configuration

| Environment Var | Description                    |
|-----------------|--------------------------------|
| `STE_API_TOKEN` | your api token                 |
| `STE_PORT`      | server port (defaults to 9119) |

The api token is a personal access token that can be created with a valid smartthings login [here](https://account.smartthings.com/tokens).

Required Oauth2 scopes
* `r:devices:*`

### Prometheus Scrape Configuration example
Since this exporter leverages the smartthings API, there is no need to target the smartthings hub directly.
```
- job_name: smartthings
  honor_timestamps: true
  scrape_interval: 15s
  scrape_timeout: 10s
  metrics_path: /metrics
  scheme: http
  follow_redirects: true
  static_configs:
  - targets:
    - smartthings-exporter:9119
```

## References

* [Smartthings Rest Api](https://smartthings.developer.samsung.com/docs/api-ref/st-api.html)