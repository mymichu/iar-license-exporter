package iardiagnostic

import (
	"log"
	"os/exec"
)

type LicenseServerInfo struct {
	Server   LicenseServer
	Licenses []License
}

type IarDiagnostics interface {
	GetIARDiagnostics() ([]LicenseServerInfo, error)
}

type iarDiagnostics struct {
	diagnosticFilePath    string
	iarlicensemanagerPath string
}

func NewIARDiagnostics(filepath string, iarlicensemanagerPath string) IarDiagnostics {
	log.Println("-- Creating new IAR Diagnostics -- ")
	log.Println(" IAR Diagnostics file path: " + filepath)
	log.Println(" IAR License Manager path: " + iarlicensemanagerPath)
	log.Println("----------------------------------")
	return &iarDiagnostics{
		diagnosticFilePath:    filepath,
		iarlicensemanagerPath: iarlicensemanagerPath}
}

func parseIARDiagnostics(filepath string) ([]LicenseServerInfo, error) {
	rawdata, error := parseFile(filepath)
	if error != nil {
		return nil, error
	}
	licenseServerInfos := make([]LicenseServerInfo, 0)
	for _, data := range rawdata.foundLicenservers[:len(rawdata.foundLicenservers)-1] {
		licenseServer, err := parseLicenseServerInfo(data.header)
		if err != nil {
			return nil, err
		}
		licenseServerInfo := LicenseServerInfo{Server: licenseServer, Licenses: make([]License, 0)}
		for _, license := range data.payload {
			licenseInfo, err := parseLicense(license)
			if err != nil {
				return nil, err
			}
			licenseServerInfo.Licenses = append(licenseServerInfo.Licenses, licenseInfo)
		}
		licenseServerInfos = append(licenseServerInfos, licenseServerInfo)
	}
	return licenseServerInfos, nil
}

func (diagnostic *iarDiagnostics) GetIARDiagnostics() ([]LicenseServerInfo, error) {

	cmd := exec.Command(diagnostic.iarlicensemanagerPath, "get_diagnostics")

	err := cmd.Run()

	if err != nil {
		return nil, err
	}
	return parseIARDiagnostics(diagnostic.diagnosticFilePath)
}
