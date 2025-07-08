package htmlx

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Parser func(string) (any, error)

// Return the raw content extracted from HTML
func StringParser(rawVal string) (any, error) {
	return rawVal, nil
}

// Return the content that has the whitespace trimmed
func StringParserClean(rawVal string) (any, error) {
	return strings.TrimSpace(rawVal), nil
}

// Return integer value of the content
func IntParser(rawVal string) (any, error) {
	trimmedRawVal := strings.TrimSpace(rawVal)
	if !regexp.MustCompile(`^[a-zA-Z$%]?\s*[0-9]+\s*[a-zA-Z$%]?$`).MatchString(trimmedRawVal) {
		return nil, fmt.Errorf("%s is not valid for parsing to integer", trimmedRawVal)
	}

	intStr := regexp.MustCompile(`[0-9]+`).FindString(trimmedRawVal)
	intVal, err := strconv.Atoi(intStr)
	if err != nil {
		return nil, err
	}

	return intVal, nil
}

// Return float value of the content
func FloatParser(rawVal string) (any, error) {
	trimmedRawVal := strings.TrimSpace(rawVal)
	if !regexp.MustCompile(`^[a-zA-Z$%]?\s*-?\d+(?:[,.]\d+)*(\.\d+)?\s*[a-zA-Z$%]?$`).
		MatchString(trimmedRawVal) {
		return nil, fmt.Errorf("%s is not valid for parsing to float", trimmedRawVal)
	}

	floatStr := regexp.MustCompile(`-?\d+(?:[,.]\d+)*(\.\d+)?`).FindString(trimmedRawVal)
	floatVal, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return nil, err
	}

	return floatVal, nil
}

// Return time.Time value of the content
func DateParser(dateFormat string) Parser {
	return func(rawVal string) (any, error) {
		date, err := time.Parse(dateFormat, strings.TrimSpace(rawVal))
		if err != nil {
			return nil, err
		}

		return date, nil
	}
}

// Set a default value if the returned content is empty, if not the alternate parser is used
func IfNullParser(defaultVal any, alternateParser Parser) Parser {
	return func(rawVal string) (any, error) {
		if strings.TrimSpace(rawVal) == "" {
			return defaultVal, nil
		}

		if alternateParser == nil {
			return nil, fmt.Errorf("The alternate parser is not provided")
		}

		return alternateParser(rawVal)
	}
}
