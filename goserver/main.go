package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	reg := prometheus.NewRegistry()
	//reg.MustRegister(collectors.NewGoCollector())
	m := NewMetrics(reg)

	m.devices.Set(float64(len(dvs)))
	m.info.With(prometheus.Labels{"version": version}).Set(1)

	dMux := http.NewServeMux()
	rdh := registerDevicesHandler{metrics: m}
	dMux.HandleFunc("/devices", rdh.registerDevices)

	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", dMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()

	select {}
}

type registerDevicesHandler struct {
	metrics *metrics
}

func (rdh *registerDevicesHandler) registerDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDevices(w, r)
	case http.MethodPost:
		createDevice(w, r, rdh.metrics)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type metrics struct {
	devices  prometheus.Gauge
	info     *prometheus.GaugeVec
	upgrades *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gomonitoring",
			Name:      "connected_devices",
			Help:      "Number of connected devices",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "gomonitoring",
			Name:      "info",
			Help:      "Information about my app environment",
		},
			[]string{"version"}),
		upgrades: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "gomonitoring",
			Name:      "device_upgrade_total",
			Help:      "Number of upgraded devices",
		}, []string{"type"}),
	}
	reg.MustRegister(m.devices, m.info, m.upgrades)
	return m
}

var dvs []*Device
var version string

func init() {
	dvs = []*Device{
		{1, "5F-22-33-44-55-66", "1.0.0"},
		{2, "5F-22-33-44-55-67", "1.0.0"},
	}
	version = "1.2.3" // usually through env
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func createDevice(w http.ResponseWriter, r *http.Request, m *metrics) {
	var dv Device

	err := json.NewDecoder(r.Body).Decode(&dv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dvs = append(dvs, &dv)
	m.devices.Set(float64(len(dvs)))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("device created"))
}

func createLogger() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)

	return log
}

func checkConnection(logger *logrus.Logger) *connectionResult {
	response, err := http.Get("https://example.com")
	if err != nil {
		logger.Debugf("http get request: %v", err)
		return &connectionResult{
			IsAlive: false,
			Message: "request error",
		}
	}
	defer response.Body.Close()

	// check status code
	if response.StatusCode != http.StatusOK {
		logger.WithFields(logrus.Fields{
			"status_code": response.StatusCode,
		}).Debug("http get request: status_code is not 200")
		return &connectionResult{
			IsAlive: false,
			Message: "example.com is not available",
		}
	}

	return &connectionResult{
		IsAlive: true,
		Message: "success",
	}
}

type connectionResult struct {
	IsAlive bool
	Message string
}

func upgradeDevice(w http.ResponseWriter, r *http.Request, m *metrics) {
	path := strings.TrimPrefix(r.URL.Path, "/devices/")

	id, err := strconv.Atoi(path)
	if err != nil || id < 1 {
		http.NotFound(w, r)
	}

	var dv Device
	err = json.NewDecoder(r.Body).Decode(&dv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range dvs {
		if dvs[]
	}
}

//func checkWebConnection(c *gin.Context) {
//	logger := createLogger()
//
//	result := checkConnection(logger)
//	if !result.IsAlive {
//		c.JSON(500, gin.H{"message": result.Message})
//		return
//	}
//	c.JSON(200, gin.H{"message": result.Message})
//}
//
//func checkWebConnectionLoop(c *gin.Context) {
//	logger := createLogger()
//
//	successCount := 0
//	for {
//		result := checkConnection(logger)
//		if !result.IsAlive {
//			c.JSON(500, gin.H{"message": result.Message})
//			return
//		}
//
//		successCount += 1
//		if successCount > 10 {
//			c.JSON(200, gin.H{"message": result.Message})
//			return
//		}
//	}
//}
