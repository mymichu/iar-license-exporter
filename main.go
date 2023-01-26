package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"iar-license-exporter/internal/iardiagnostic"
	"iar-license-exporter/internal/localhostserver"
	prometheusmetrics "iar-license-exporter/internal/prometheus-metrics"
)

func main() {

	port := flag.String("port", "9101", "Port to listen on")
	diagnosticFilepath := flag.String("diagnostic-file", "/usr/local/etc/IARSystems/LicenseManagement/Diagnostics/Diagnostics.txt", "Path to diagnostic file")
	iarManagerLicensePath := flag.String("lms-path", "iar_licenses_used", "Path to iar license manager binary")

	flag.Parse()

	serverSpecs := localhostserver.GetServerSpecs()
	iardiagnostic := iardiagnostic.NewIARDiagnostics(*diagnosticFilepath, *iarManagerLicensePath)

	prometheusmetrics.NewIARMetrics(iardiagnostic, prometheusmetrics.ServerInfo{Hostname: serverSpecs.Hostname, Ipaddr: serverSpecs.Ipaddr})
	log.Println("Starting server on port: ", *port)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
