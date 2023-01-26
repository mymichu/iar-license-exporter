package iardiagnostic

import (
	"errors"
	"strings"
)

type LicenseServer struct {
	Version  string
	IP       string
	Hostname string
}

var keywordsLicensServer = map[string]func(text string, licenseServer *LicenseServer) error{
	"Version:": func(text string, licenseServer *LicenseServer) error {
		licenseServer.Version = text
		return nil
	},
	"IP/Host name:": func(text string, licenseServer *LicenseServer) error {
		ip_hostName := strings.Split(text, "/")
		if len(ip_hostName) != 2 {
			return errors.New("IP/Host name is not in the correct format")
		}
		licenseServer.IP = ip_hostName[0]
		licenseServer.Hostname = ip_hostName[1]
		return nil
	},
}

func parseLicenseServerInfo(text string) (LicenseServer, error) {
	if text == "" {
		return LicenseServer{}, errors.New("License Server to analyze is empty")
	}
	var licenseServer LicenseServer

	for k, action := range keywordsLicensServer {
		firstIndex := strings.Index(text, k) + len(k)
		lastIndex := strings.Index(text[firstIndex:], "[")
		if firstIndex > len(k) {
			valueOfInterest := text[firstIndex:]
			if lastIndex != -1 {
				valueOfInterest = text[firstIndex : firstIndex+lastIndex]
			}
			valueOfInterest = strings.TrimSpace(valueOfInterest)
			err := action(valueOfInterest, &licenseServer)
			if err != nil {
				return licenseServer, err
			}
		}
	}
	return licenseServer, nil
}
