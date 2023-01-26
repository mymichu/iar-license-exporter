package iardiagnostic

import (
	"bufio"
	"os"
	"strings"
)

const (
	None = iota
	SystemInfo
	NewPayload
	Payload
	NewLicenserverHeader
	LicenserverHeader
	InstalledProducts
	LicenseRequestTests
	CommuterLicense
)

func getNextState(line string, currentState int) int {
	if currentState == NewPayload {
		return Payload
	}
	if currentState == NewLicenserverHeader {
		return LicenserverHeader
	}
	if strings.Contains(line, " #################################################################") {
		switch currentState {
		case SystemInfo:
			return InstalledProducts
		case InstalledProducts:
			return LicenseRequestTests
		case LicenseRequestTests:
			return CommuterLicense
		case CommuterLicense:
			return NewLicenserverHeader
		case Payload:
			return NewLicenserverHeader
		case None:
			return SystemInfo
		default:
			return currentState
		}
	} else if strings.Contains(line, " --------------------") {
		return NewPayload
	}
	return currentState
}

type serverDataRaw struct {
	header  string
	payload []string
}
type rawData struct {
	systemInfo        string
	installedProducts string
	licenseRequest    string
	commuterLicense   string
	foundLicenservers []serverDataRaw
}

func parseFile(filepath string) (rawData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return rawData{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	servers := make([]serverDataRaw, 0)
	systemInfo := ""
	installedProducts := ""
	licenseRequestTests := ""
	commuterLicense := ""
	currentState := None
	for scanner.Scan() {
		line := scanner.Text()
		currentState = getNextState(line, currentState)
		switch currentState {
		case SystemInfo:
			systemInfo += line
		case NewLicenserverHeader:
			servers = append(servers, serverDataRaw{header: "", payload: make([]string, 0)})
		case LicenserverHeader:
			servers[len(servers)-1].header += line
		case Payload:
			servers[len(servers)-1].payload[len(servers[len(servers)-1].payload)-1] += line
		case NewPayload:
			servers[len(servers)-1].payload = append(servers[len(servers)-1].payload, "")
		case InstalledProducts:
			installedProducts += line
		case LicenseRequestTests:
			licenseRequestTests += line
		case CommuterLicense:
			commuterLicense += line
		}
	}
	if err := scanner.Err(); err != nil {
		return rawData{}, err
	}
	return rawData{systemInfo: systemInfo,
		installedProducts: installedProducts,
		licenseRequest:    licenseRequestTests,
		commuterLicense:   commuterLicense,
		foundLicenservers: servers}, nil
}
