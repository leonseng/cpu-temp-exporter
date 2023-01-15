package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 80, "port to listen on")
	flag.Parse()

	_, err := exec.LookPath("sensors")
	if err != nil {
		fmt.Println("lm_sensors package is not installed.")
		os.Exit(1)
	}

	// define a regular expression to match the temperature values
	re := regexp.MustCompile(`\+[0-9]{2}\.[0-9]{1}°C`)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		// execute the `sensors` command and capture its output
		cmdOutput, err := exec.Command("sensors").Output()
		if err != nil {
			fmt.Println("An error occurred while running the sensors command:", err)
			fmt.Fprintf(w, "cpu_temp_sensors_up 0\n")
			return
		}

		// initialize a map to store the CPU temperatures
		cpuTemps := make(map[string]float64)

		// split the output by newline
		for _, line := range strings.Split(string(cmdOutput), "\n") {
			if strings.Contains(line, "Core") {
				// extract the temperature value from lines for CPU cores
				tempStr := re.FindString(line)
				tempStr = strings.TrimPrefix(tempStr, "+")
				tempStr = strings.TrimSuffix(tempStr, "°C")
				temp, err := strconv.ParseFloat(tempStr, 64)
				if err != nil {
					fmt.Println("An error occurred while converting the temperature value to a float:", err)
					fmt.Fprintf(w, "cpu_temp_sensors_up 0\n")
					return
				}

				// store temperature value in map indexed by core number
				coreNum := strings.TrimPrefix(strings.Split(line, ":")[0], "Core ")
				cpuTemps[coreNum] = temp
			}
		}

		// retrieve hostname to be added as label. Defined within this handler as hostname could change
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println("An error occurred while getting the hostname:", err)
			fmt.Fprintf(w, "sensors_up 0\n")
			return
		}

		// write the Prometheus metrics to the HTTP response
		fmt.Fprintf(w, "cpu_temp_sensors_up 1\n")
		for core, temp := range cpuTemps {
			metric := fmt.Sprintf("%.1f", temp)
			fmt.Fprintf(w, "cpu_temp{core=\"%s\", hostname=\"%s\"} %s\n", core, hostname, metric)
		}
	})

	// start the HTTP server
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
