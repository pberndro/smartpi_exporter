package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/pberndro/smartpi_exporter/ade7878"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/io/i2c"
)


func makeReadoutAccumulator() (r ReadoutAccumulator) {
	r.Current = make(ade7878.Readings)
	r.Voltage = make(ade7878.Readings)
	r.ActiveWatts = make(ade7878.Readings)
	r.CosPhi = make(ade7878.Readings)
	r.Frequency = make(ade7878.Readings)
	r.WattHoursConsumed = make(ade7878.Readings)
	r.WattHoursProduced = make(ade7878.Readings)
	return r
}

func makeReadout() (r ade7878.ADE7878Readout) {
	r.Current = make(ade7878.Readings)
	r.Voltage = make(ade7878.Readings)
	r.ActiveWatts = make(ade7878.Readings)
	r.CosPhi = make(ade7878.Readings)
	r.Frequency = make(ade7878.Readings)
	r.ApparentPower = make(ade7878.Readings)
	r.ReactivePower = make(ade7878.Readings)
	r.PowerFactor = make(ade7878.Readings)
	r.ActiveEnergy = make(ade7878.Readings)
	return r
}

func pollSmartPi(config *ade7878.Config, device *i2c.Device) {
	var wattHourBalanced float64
	var p ade7878.Phase

	accumulator := makeReadoutAccumulator()
	i := 0

	tick := time.Tick(time.Duration(1000/config.Samplerate) * time.Millisecond)

	for {
		readouts := makeReadout()
		// Restart the accumulator loop every 60 seconds.
		if i > (60*config.Samplerate - 1) {
			i = 0
			accumulator = makeReadoutAccumulator()
		}

		startTime := time.Now()

		// Update readouts and the accumlator.
		ade7878.ReadPhase(device, config, ade7878.PhaseN, &readouts)
		accumulator.Current[ade7878.PhaseN] += readouts.Current[ade7878.PhaseN] / (60.0 * float64(config.Samplerate))

		for _, p = range ade7878.MainPhases {
			ade7878.ReadPhase(device, config, p, &readouts)
			accumulator.Current[p] += readouts.Current[p] / (60.0 * float64(config.Samplerate))
			accumulator.Voltage[p] += readouts.Voltage[p] / (60.0 * float64(config.Samplerate))
			accumulator.ActiveWatts[p] += readouts.ActiveWatts[p] / (60.0 * float64(config.Samplerate))
			accumulator.CosPhi[p] += readouts.CosPhi[p] / (60.0 * float64(config.Samplerate))
			accumulator.Frequency[p] += readouts.Frequency[p] / (60.0 * float64(config.Samplerate))

			if readouts.ActiveWatts[p] >= 0 {
				accumulator.WattHoursConsumed[p] += math.Abs(readouts.ActiveWatts[p]) / (3600.0 * float64(config.Samplerate))
			} else {
				accumulator.WattHoursProduced[p] += math.Abs(readouts.ActiveWatts[p]) / (3600.0 * float64(config.Samplerate))
			}
			wattHourBalanced += readouts.ActiveWatts[p] / (3600.0 * float64(config.Samplerate))
		}

		// Update metrics endpoint.
		updatePrometheusMetrics(&readouts)

		/*
		// Every sample
		if i%1 == 0 {
			if config.SharedFileEnabled {
				writeSharedFile(config, &readouts, wattHourBalanced)
			}
			wattHourBalanced = 0
		}
		*/

		// Every 60 seconds.
		if i == (60*config.Samplerate - 1) {

			// balanced value
			var wattHourBalanced60s float64
			//consumedWattHourBalanced60s = 0.0
			//producedWattHourBalanced60s = 0.0

			for _, p = range ade7878.MainPhases {
				wattHourBalanced60s += accumulator.WattHoursConsumed[p]
				wattHourBalanced60s -= accumulator.WattHoursProduced[p]
			}
			//if wattHourBalanced60s >= 0 {
				//consumedWattHourBalanced60s = wattHourBalanced60s
			//} else {
				//producedWattHourBalanced60s = wattHourBalanced60s
			//}
		}

		delay := time.Since(startTime) - (time.Duration(1000/config.Samplerate) * time.Millisecond)
		if int64(delay) > 0 {
			log.Errorf("Readout delayed: %s", delay)
		}
		<-tick
		i++
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	prometheus.MustRegister(currentMetric)
	prometheus.MustRegister(voltageMetric)
	prometheus.MustRegister(activePowerMetirc)
	prometheus.MustRegister(cosphiMetric)
	prometheus.MustRegister(frequencyMetric)
	prometheus.MustRegister(apparentPowerMetric)
	prometheus.MustRegister(reactivePowerMetric)
	prometheus.MustRegister(powerFactorMetric)
	prometheus.MustRegister(version.NewCollector("smartpi"))
}

var appVersion = "No Version Provided"

func main() {
	config := ade7878.NewConfig()

	version := flag.Bool("v", false, "prints current version information")
	flag.Parse()
	if *version {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	listenAddress := config.MetricsListenAddress
	device, _ := ade7878.InitADE7878(config)

	go pollSmartPi(config, device)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
            <head><title>SmartPi Readout Metrics Server</title></head>
            <body>
            <h1>SmartPi Readout Metrics Server</h1>
            <p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		panic(fmt.Errorf("Error starting HTTP server: %s", err))
	}
}
