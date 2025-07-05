package htmlx

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Parser func(string) (any, error)

type Config struct {
	dateFormat string
	parsers    map[string]Parser
}

func NewDefaultConfig() *Config {
	return &Config{
		dateFormat: "2006-01-02T15:04:05Z07:00",
		parsers:    map[string]Parser{},
	}
}

type Option func(*Config)

// Set the date format which is used for time.Time field, the default one is "2006-01-02T15:04:05Z07:00"
func SetDateFormat(format string) Option {
	return func(c *Config) {
		c.dateFormat = format
	}
}

func SetParsers(parsers map[string]Parser) Option {
	return func(c *Config) {
		maps.Copy(c.parsers, parsers)
	}
}

type HtmlxTags struct {
	selector string
	source   string
	parser   string
}

func initializeHtmlxTags(fieldType reflect.StructField) (HtmlxTags, error) {
	var htmlxTags HtmlxTags

	htmlxTags.selector = fieldType.Tag.Get("selector")
	if htmlxTags.selector == "" {
		return htmlxTags, fmt.Errorf("Missing selector for field '%s'", fieldType.Name)
	}

	htmlxTags.source = fieldType.Tag.Get("source")
	if htmlxTags.source == "" {
		htmlxTags.source = "content"
	}

	htmlxTags.parser = fieldType.Tag.Get("parser")

	return htmlxTags, nil
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
		// TODO: Reigister all struct tags here
		htmlxTags, err := initializeHtmlxTags(fieldType)
		if err != nil {
			return fmt.Errorf("Error extracting tags from field '%s': %s", fieldType.Name, err.Error())
		}

		htmlElement := sel.Find(htmlxTags.selector)
		if htmlElement.Length() == 0 {
			return fmt.Errorf("Error locating html element for field '%s'", fieldType.Name)
		}

		rawVal, err := getRawValue(fieldType, htmlElement, htmlxTags.source)
		if err != nil {
			return fmt.Errorf("Error getting raw value from field '%s': %s", fieldType.Name, err.Error())
		}

		if err = parseValue(fieldVal, rawVal, config, htmlxTags); err != nil {
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

func stringParser(fieldVal reflect.Value, rawVal string) {
	fieldVal.SetString(strings.TrimSpace(rawVal))
}

func intParser(fieldVal reflect.Value, rawVal string) error {
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

	return nil
}

func floatParser(fieldVal reflect.Value, rawVal string) error {
	trimmedRawVal := strings.TrimSpace(rawVal)
	if !regexp.MustCompile(`^[a-zA-Z$%]?\s*-?\d+(?:[,.]\d+)*(\.\d+)?\s*[a-zA-Z$%]$`).
		MatchString(trimmedRawVal) {
		return fmt.Errorf("%s is not valid for parsing to float", trimmedRawVal)
	}

	floatStr := regexp.MustCompile(`-?\d+(?:[,.]\d+)*(\.\d+)?`).FindString(trimmedRawVal)
	floatVal, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return err
	}

	fieldVal.SetFloat(floatVal)

	return nil
}

func dateParser(fieldVal reflect.Value, rawVal, dateFormat string) error {
	date, err := time.Parse(dateFormat, strings.TrimSpace(rawVal))
	if err != nil {
		return err
	}

	fieldVal.Set(reflect.ValueOf(date))

	return nil
}

func parseValue(fieldVal reflect.Value, rawVal string, config *Config, htmlxTags HtmlxTags) error {
	var err error

	if !fieldVal.IsValid() {
		return fmt.Errorf("Field doesn't represent a value")
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Field can't be set")
	}

	if htmlxTags.parser != "" {
		parser, ok := config.parsers[htmlxTags.parser]
		if !ok {
			return fmt.Errorf("parser %s is not recognizable", htmlxTags.parser)
		}

		val, err := parser(rawVal)
		if err != nil {
			return fmt.Errorf("parser '%s' error: %s", htmlxTags.parser, err.Error())
		}

		processedVal := reflect.ValueOf(val)
		if !processedVal.IsValid() {
			return fmt.Errorf("processed value using parser %s is invalid", htmlxTags.parser)
		}

		if !processedVal.Type().AssignableTo(fieldVal.Type()) {
			return fmt.Errorf(
				"Incompatible type when using parser %s, want %s, get %s",
				htmlxTags.parser,
				fieldVal.Type().String(),
				processedVal.Type().String(),
			)
		}

		fieldVal.Set(processedVal)

		return nil
	}

	switch fieldVal.Type() {
	case reflect.TypeOf(""):
		stringParser(fieldVal, rawVal)
	case reflect.TypeOf(0):
		if err = intParser(fieldVal, rawVal); err != nil {
			return err
		}
	case reflect.TypeOf(0.5):
		if err = floatParser(fieldVal, rawVal); err != nil {
			return err
		}
	case reflect.TypeOf(time.Time{}):
		if err = dateParser(fieldVal, rawVal, config.dateFormat); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Value of type %s is not supported", fieldVal.Type().String())
	}

	return nil
}
