package prometheusmetrics

import (
	"log"
	"time"

	"iar-license-exporter/internal/iardiagnostic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ServerInfo struct {
	Hostname string
	Ipaddr   string
}

type iarMetrics struct {
	diagnostics iardiagnostic.IarDiagnostics
}

func (metrics *iarMetrics) trackLicenses(used *prometheus.GaugeVec, available *prometheus.GaugeVec, total *prometheus.GaugeVec, inQueue *prometheus.GaugeVec, reserved *prometheus.GaugeVec) {
	go func() {
		for {
			serverInfos, error := metrics.diagnostics.GetIARDiagnostics()
			if error != nil {
				log.Fatal(error)
			}
			for _, serverInfo := range serverInfos {
				hostname := serverInfo.Server.Hostname
				ipaddr := serverInfo.Server.IP
				for _, licenseInfo := range serverInfo.Licenses {
					used.WithLabelValues(hostname, ipaddr, licenseInfo.Feature, licenseInfo.LicenseNumber, licenseInfo.Version).Set(float64(licenseInfo.NumberTaken))
					available.WithLabelValues(hostname, ipaddr, licenseInfo.Feature, licenseInfo.LicenseNumber, licenseInfo.Version).Set(float64(licenseInfo.AvailableLicenses))
					total.WithLabelValues(hostname, ipaddr, licenseInfo.Feature, licenseInfo.LicenseNumber, licenseInfo.Version).Set(float64(licenseInfo.NumberOfLicenses))
					inQueue.WithLabelValues(hostname, ipaddr, licenseInfo.Feature, licenseInfo.LicenseNumber, licenseInfo.Version).Set(float64(licenseInfo.LicensesInQueue))
					reserved.WithLabelValues(hostname, ipaddr, licenseInfo.Feature, licenseInfo.LicenseNumber, licenseInfo.Version).Set(float64(licenseInfo.Reserved))
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func (metrics *iarMetrics) registerLicenseMetrics(host ServerInfo) {

	licenseUsed := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "licenses_used",
		Help:        "The total number of licenses used",
		ConstLabels: map[string]string{"instance": host.Hostname, "ipaddr": host.Ipaddr},
	}, []string{"license_server_hostname", "license_server_ipaddr", "feature", "license_number", "version"})

	licensesAvailable := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "licenses_available",
		Help:        "The total number of licenses used",
		ConstLabels: map[string]string{"instance": host.Hostname, "ipaddr": host.Ipaddr},
	}, []string{"license_server_hostname", "license_server_ipaddr", "feature", "license_number", "version"})

	licensesTotalNumber := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "licenses_total_number",
		Help:        "The total number of licenses used",
		ConstLabels: map[string]string{"instance": host.Hostname, "ipaddr": host.Ipaddr},
	}, []string{"license_server_hostname", "license_server_ipaddr", "feature", "license_number", "version"})

	licensesInQueue := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "licenses_in_queue",
		Help:        "The total number of licenses used",
		ConstLabels: map[string]string{"instance": host.Hostname, "ipaddr": host.Ipaddr},
	}, []string{"license_server_hostname", "license_server_ipaddr", "feature", "license_number", "version"})

	licensesReserved := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "licenses_reserved",
		Help:        "The total number of licenses used",
		ConstLabels: map[string]string{"instance": host.Hostname, "ipaddr": host.Ipaddr},
	}, []string{"license_server_hostname", "license_server_ipaddr", "feature", "license_number", "version"})

	metrics.trackLicenses(licenseUsed, licensesAvailable, licensesTotalNumber, licensesInQueue, licensesReserved)

}

func NewIARMetrics(diagnostics iardiagnostic.IarDiagnostics, hostServer ServerInfo) {
	metrics := iarMetrics{diagnostics: diagnostics}
	metrics.registerLicenseMetrics(hostServer)
}
