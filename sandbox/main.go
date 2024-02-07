package main

import (
	"fmt"
	"strings"
	"unicode"
)

func deriveLayoutFromDate(dateString string) (string, error) {
	// Define a mapping of specifiers to their corresponding values
	specifiers := map[string]string{
		"2006":      "%Y",
		"06":        "%y",
		"01":        "%m",
		"january":   "%B",
		"jan":       "%b",
		"02":        "%d",
		"1":         "%u",
		"monday":    "%A",
		"mon":       "%a",
		"tuesday":   "%A",
		"tue":       "%a",
		"wednesday": "%A",
		"wed":       "%a",
		"thursday":  "%A",
		"thu":       "%a",
		"friday":    "%A",
		"fri":       "%a",
		"saturday":  "%A",
		"sat":       "%a",
		"sunday":    "%A",
		"sun":       "%a",
		"15":        "%H",
		"03":        "%I",
		"04":        "%M",
		"05":        "%S",
		"3":         "%3f",
		"6":         "%6f",
		"9":         "%9f",
		"pm":        "%P",
		"am":        "%P",
		"PM":        "%p",
		"AM":        "%p",
		"-0700":     "%z",
		"-07:00":    "%Ez",
		"ce":        "%EC",
	}

	// Split the input date string into date and time components
	dateTimeParts := strings.FieldsFunc(dateString, func(r rune) bool {
		return !unicode.IsDigit(r) && !unicode.IsLetter(r)
	})

	// Derive the layout based on the input date string
	layout := ""
	for _, part := range dateTimeParts {
		lowerPart := strings.ToLower(part)
		if specifier, ok := specifiers[lowerPart]; ok {
			layout += specifier
		} else if unicode.IsDigit([]rune(part)[0]) {
			// If the part is a numeric value, use it as is
			layout += part
		} else {
			// If the part is not found in specifiers and not a numeric value,
			// assume it's a valid component without raising an error
			layout += part
		}
	}

	// Additional notes: %j, %G, %V
	layout = strings.ReplaceAll(layout, "%j", "002") // Note 2
	layout = strings.ReplaceAll(layout, "%G", "%Y")  // Note 2
	layout = strings.ReplaceAll(layout, "%V", "01")  // Note 2

	return layout, nil
}

func main() {
	// Example usage
	inputDateString := "2020"
	result, err := deriveLayoutFromDate(inputDateString)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Derived Layout:", result)
	}
}
