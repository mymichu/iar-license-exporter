package iardiagnostic

import (
	"errors"
	"strconv"
	"strings"
)

type License struct {
	Feature                 string
	Version                 string
	NumberOfLicenses        int
	NumberTaken             int
	IsCommuter              bool
	MaxCommuterCheckoutDays int
	LicensesInQueue         int
	Reserved                int
	AvailableLicenses       int
	LicenseNumber           string
}

func convertToInt(text string, errorMessage string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		return -1, errors.New(errorMessage)
	}
	return number, nil
}

func parseFeature(text string, license *License) error {

	license.Feature = text
	return nil
}

func parseVersion(text string, license *License) error {
	license.Version = text
	return nil
}

func parseNumberOfLicenses(text string, license *License) error {
	numberLicenses, err := convertToInt(text, "Number of licenses is not a number")
	if err != nil {
		return err
	}
	license.NumberOfLicenses = numberLicenses
	return nil
}

func parseNumberOfTake(text string, license *License) error {
	numberTake, err := convertToInt(text, "Number of licenses taken is not a number")
	if err != nil {
		return err
	}
	license.NumberTaken = numberTake
	return nil
}

func parseIsCommuter(text string, license *License) error {
	license.IsCommuter = text == "yes"
	return nil
}

func parseMaxCommuterCheckout(text string, license *License) error {
	maxCheckoutDays, err := convertToInt(text, "Max commuter checkout days is not a number")
	if err != nil {
		return err
	}
	license.MaxCommuterCheckoutDays = maxCheckoutDays
	return nil
}

func parseLicensesInQueue(text string, license *License) error {
	licensesInQueue, err := convertToInt(text, "Licenses in queue is not a number")
	if err != nil {
		return err
	}
	license.LicensesInQueue = licensesInQueue
	return nil
}

func numberLicensesReserved(text string, license *License) error {
	numberLicensesReserved, err := convertToInt(text, "Number of licenses reserved is not a number")
	if err != nil {
		return err
	}
	license.Reserved = numberLicensesReserved
	return nil
}

func parseAvailableLicenses(text string, license *License) error {
	numberAvailableLicense, err := convertToInt(text, "Number of available licenses is not a number")
	if err != nil {
		return err
	}
	license.AvailableLicenses = numberAvailableLicense
	return nil
}

func parseLicenseNumber(text string, license *License) error {
	license.LicenseNumber = text
	return nil
}

var keywords = map[string]func(text string, license *License) error{
	"Feature:":                 parseFeature,
	"Version:":                 parseVersion,
	"Numberoflicenses:":        parseNumberOfLicenses,
	"Numbertaken:":             parseNumberOfTake,
	"Iscommuter:":              parseIsCommuter,
	"Maxcommutercheckoutdays:": parseMaxCommuterCheckout,
	"Licensesinqueue:":         parseLicensesInQueue,
	"Reserved:":                numberLicensesReserved,
	"Availablelicenses:":       parseAvailableLicenses,
	"Licensenumber:":           parseLicenseNumber,
}

func parseLicense(text string) (License, error) {
	text = strings.ReplaceAll(text, " ", "")
	var license License

	for k, action := range keywords {
		firstIndex := strings.Index(text, k) + len(k)
		lastIndex := strings.Index(text[firstIndex:], "[")
		if firstIndex > len(k) {
			valueOfInterest := text[firstIndex:]
			if lastIndex != -1 {
				valueOfInterest = text[firstIndex : firstIndex+lastIndex]
			}
			err := action(valueOfInterest, &license)
			if err != nil {
				return license, err
			}
		}
	}
	return license, nil
}
