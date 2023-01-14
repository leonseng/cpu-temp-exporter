package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	_, err := exec.LookPath("sensors")
	if err != nil {
		fmt.Println("lm_sensors package is not installed.")
		os.Exit(1)
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// execute the `sensors` command and capture its output
		output, _ := exec.Command("sensors").Output()
		// convert the output to a string
		outputStr := string(output)

		// initialize a map to store the CPU temperatures
		cpuTemps := make(map[string]float64)

		// define a regular expression to match the temperature values
		re := regexp.MustCompile(`\+[0-9]{2}\.[0-9]{1}°C`)

		// split the output by newline
		for _, line := range strings.Split(outputStr, "\n") {
			// check if the line contains "Core"
			if strings.Contains(line, "Core") {
				// extract the temperature value from the line
				tempStr := re.FindString(line)
				// remove the '+' and '°C' characters from the temperature value
				tempStr = strings.TrimPrefix(tempStr, "+")
				tempStr = strings.TrimSuffix(tempStr, "°C")
				// convert the temperature value to a float
				temp, _ := strconv.ParseFloat(tempStr, 64)
				// extract the core number from the line
				coreNum := strings.TrimPrefix(strings.Split(line, ":")[0], "Core ")
				// add the temperature value to the map
				cpuTemps[coreNum] = temp
			}
		}
		hostname, _ := os.Hostname()
		// write the Prometheus metrics to the HTTP response
		for core, temp := range cpuTemps {
			metric := fmt.Sprintf("%.1f", temp)
			fmt.Fprintf(w, `cpu_temp{core="%s", hostname="%s"} %s\n`, core, hostname, metric)
		}
	})

	// start the HTTP server
	addr := ":9090"
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
