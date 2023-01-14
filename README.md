# CPU Temperature Exporter

This Go code exposes an HTTP endpoint that can be scraped by Prometheus for CPU temperature metrics.

When the service is running, Prometheus will be able to scrape the metrics by visiting http://localhost:9090/metrics or the hostname or IP of the machine where the service is running instead of localhost.

Here is an example of the metrics that will be returned:
```
$ curl -s http://localhost:9090/metrics
cpu_temp_sensors_up 1
cpu_temp{core="0", hostname="hostname"} 60.0
cpu_temp{core="1", hostname="hostname"} 60.0
cpu_temp{core="2", hostname="hostname"} 60.0
cpu_temp{core="3", hostname="hostname"} 59.0
```

It runs the `sensors` command from the [lm_sensors](https://wiki.archlinux.org/title/lm_sensors) package on the local machine, extracts the CPU temperatures from its output, and returns them as Prometheus metrics.

## Dependencies

This code requires the following dependencies:

- Go 1.15+
- [lm_sensors](https://wiki.archlinux.org/title/lm_sensors) package

## How to Build

1. Clone the repository and navigate to the project folder:
    ```
    git clone https://github.com/leonseng/cpu-temperature-exporter.git
    cd cpu-temperature-exporter
    ```
1. Build the binary:
    ```
    GOOS=linux GOARCH=amd64 go build -o cpu_temp_exporter
    ```
    This will produce a binary named cpu_temp_exporter that runs on Linux platform with AMD64 architecture.

    You can also use the GOOS and GOARCH environment variables to specify the target operating system and architecture respectively. You can check the available options of GOARCH and GOOS by running the command `go tool dist list`.

## How to run

> Make sure you have lm-sensors package installed on your machine before running the service.

Run the binary:
```
./cpu_temp_exporter
```

The service will start running on port 9090, and Prometheus will be able to scrape the metrics by visiting http://localhost:9090/metrics or the hostname or IP of the machine where the service is running instead of localhost.
