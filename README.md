# URL shortener

-   This application is able to shorten long url to short url by encoding the generated random uuid-v4 to base64
-   This application is also used for demonstrating the application instrumentation, especially structured logging, trace, and metric using open telemetry.

## Running it locally

### Run Docker Compose

    ```shell
    docker-compose up
    ```

### Open Dashboard

1. Jaeger Tracing Dashboard: [http://127.0.0.1:16686/](http://127.0.0.1:16686/)
2. Prometheus Dashboard: [http://127.0.0.1:9090/](http://127.0.0.1:9090/)
3. Grafana Dashboard: [http://127.0.0.1:3000/](http://127.0.0.1:3000/)

## Dependencies

-   net/http
-   google/uuid
-   opentelemetry
-   rs/zerolog
-   spf13/viper
-   sqlc-dev/sqlc
