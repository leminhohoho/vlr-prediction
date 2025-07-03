package htmlx

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	dateFormat string
}

func NewDefaultConfig() *Config {
	return &Config{
		dateFormat: "2006-01-02T15:04:05Z07:00",
	}
}

type Option func(*Config)

// Set the date format which is used for time.Time field, the default one is "2006-01-02T15:04:05Z07:00"
func SetDateFormat(format string) Option {
	return func(c *Config) {
		c.dateFormat = format
	}
}

// Parse the content from HTML string to struct
func ParseFromString(
	s any,
	content string,
	opts ...Option,
) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return err
	}
	return ParseFromSelection(s, doc.Selection, opts...)
}

// Parse the content from goquery Selection to struct
func ParseFromDocument(
	s any,
	doc *goquery.Document,
	opts ...Option,
) error {
	return ParseFromSelection(s, doc.Selection, opts...)
}

// Parse the content from HTML string to struct with context
func ParseFromStringWithContext(
	ctx context.Context,
	s any,
	content string,
	opts ...Option,
) error {
	errChan := make(chan error)

	go func() {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			errChan <- err
		}

		errChan <- ParseFromSelection(s, doc.Selection, opts...)
	}()

	select {
	case err := <-errChan:
		return err
	}
}

// Parse the content from goquery Document to struct with context
func ParseFromDocumentWithContext(
	ctx context.Context,
	s any,
	doc *goquery.Document,
	opts ...Option,
) error {
	errChan := make(chan error)

	go func() {
		errChan <- ParseFromSelection(s, doc.Selection, opts...)
	}()

	select {
	case err := <-errChan:
		return err
	}
}

// Parse the content from goquery Selection to struct with context
func ParseFromSelectionWithContext(
	ctx context.Context,
	s any,
	sel *goquery.Selection,
	opts ...Option,
) error {
	errChan := make(chan error)

	go func() {
		errChan <- ParseFromSelection(s, sel, opts...)
	}()

	select {
	case err := <-errChan:
		return err
	}
}

// Parse the content from goquery Document to struct
func ParseFromSelection(
	s any,
	sel *goquery.Selection,
	opts ...Option,
) error {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	v = v.Elem()

	config := NewDefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	if err := parseFromReflectValue(v, sel, config); err != nil {
		return err
	}

	return nil
}

// Check if the struct is of supported struct types
func isStructToParse(v reflect.Value) bool {
	switch v.Type() {
	case reflect.TypeOf(time.Time{}):
		return true
		// NOTE: Add supported types here
	}

	return false
}

func parseFromReflectValue(
	v reflect.Value,
	sel *goquery.Selection,
	config *Config,
) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("value is not a struct")
	}
	t := v.Type()

	for i := range t.NumField() {
		var err error

		fieldType := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct && !isStructToParse(fieldVal) {
			if err = parseFromReflectValue(fieldVal, sel, config); err != nil {
				return fmt.Errorf("Error parsing value to field '%s' : %s", fieldType.Name, err.Error())
			}
			continue
		}

		selector := fieldType.Tag.Get("selector")
		if selector == "" {
			return fmt.Errorf("Missing selector for field '%s'", fieldType.Name)
		}

		htmlElement := sel.Find(selector)
		if htmlElement.Length() == 0 {
			return fmt.Errorf("Error locating html element for field '%s'", fieldType.Name)
		}

		source := fieldType.Tag.Get("source")
		if source == "" {
			source = "content"
		}

		rawVal, err := getRawValue(fieldType, htmlElement, source)
		if err != nil {
			return fmt.Errorf("Error getting raw value from field '%s': %s", fieldType.Name, err.Error())
		}

		if err = parseValue(&fieldVal, rawVal, config); err != nil {
			return fmt.Errorf("Error parsing value to field '%s': %s", fieldType.Name, err.Error())
		}
	}

	return nil
}

func getRawValue(
	fieldType reflect.StructField,
	htmlElement *goquery.Selection,
	source string,
) (string, error) {
	if source == "content" {
		return htmlElement.Children().Remove().End().Text(), nil
	} else if regexp.MustCompile(`^attr=[a-zA-Z-0-9]+$`).MatchString(source) {
		var exists bool
		attrName := source[5:]
		value, exists := htmlElement.Attr(attrName)
		if !exists {
			return "", fmt.Errorf("Error locating attribute %s for field '%s'", source, fieldType.Name)
		}

		return value, nil
	} else {
		return "", fmt.Errorf("unrecognizable source %s for field '%s'", source, fieldType.Name)
	}
}

func parseValue(fieldVal *reflect.Value, rawVal string, config *Config) error {
	if !fieldVal.IsValid() {
		return fmt.Errorf("Field doesn't represent a value")
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Field can't be set")
	}

	switch fieldVal.Type() {
	case reflect.TypeOf(""):
		fieldVal.SetString(strings.TrimSpace(rawVal))
	case reflect.TypeOf(0):
		trimmedRawVal := strings.TrimSpace(rawVal)
		if !regexp.MustCompile(`^[a-zA-Z$%]?\s*[0-9]+\s*[a-zA-Z$%]?`).MatchString(trimmedRawVal) {
			return fmt.Errorf("%s is not valid for parsing to integer", trimmedRawVal)
		}

		intStr := regexp.MustCompile(`[0-9]+`).FindString(trimmedRawVal)
		intVal, err := strconv.Atoi(intStr)
		if err != nil {
			return err
		}

		fieldVal.SetInt(int64(intVal))
	case reflect.TypeOf(0.5):
		trimmedRawVal := strings.TrimSpace(rawVal)
		if !regexp.MustCompile(`^[a-zA-Z$%]?\s*-?\d+(?:[,.]\d+)*(\.\d+)?\s*[a-zA-Z$%]$`).
			MatchString(trimmedRawVal) {
			return fmt.Errorf("%s is not valid for parsing to integer", trimmedRawVal)
		}

		floatStr := regexp.MustCompile(`-?\d+(?:[,.]\d+)*(\.\d+)?`).FindString(trimmedRawVal)
		floatVal, err := strconv.ParseFloat(floatStr, 64)
		if err != nil {
			return err
		}

		fieldVal.SetFloat(floatVal)
	case reflect.TypeOf(time.Time{}):
		date, err := time.Parse(config.dateFormat, strings.TrimSpace(rawVal))
		if err != nil {
			return err
		}

		fieldVal.Set(reflect.ValueOf(date))
	default:
		return fmt.Errorf("Value of type %s is not supported", fieldVal.Type().String())
	}

	return nil
}
