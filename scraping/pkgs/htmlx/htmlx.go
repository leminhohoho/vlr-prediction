package htmlx

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	dateFormat          string
	parsers             map[string]Parser
	noEmptySelection    bool
	parseAllFields      bool
	noMissingAttributes bool
	noPassThroughStruct bool
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

// Set custom parsers for values, name of the parser is corresponded o the name in parser field of the struct tags
func SetParsers(parsers map[string]Parser) Option {
	return func(c *Config) {
		maps.Copy(c.parsers, parsers)
	}
}

// Set the parser o throw error if a selector return empty selection
func SetNoEmptySelection(forbid bool) Option {
	return func(c *Config) {
		c.noEmptySelection = forbid
	}
}

// Set the parser to force all fields to be parsable (not having required struct tags will throw error)
func SetParseAllFields(strict bool) Option {
	return func(c *Config) {
		c.parseAllFields = strict
	}
}

// Force the attributes if specified as sources to must exist
func SetNoMissingAttributes(forbid bool) Option {
	return func(c *Config) {
		c.noMissingAttributes = forbid
	}
}

// Forbid the parser from accessing inner structs
func SetNoPassThroughStruct(forbid bool) Option {
	return func(c *Config) {
		c.noPassThroughStruct = forbid
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
			if !config.noPassThroughStruct {
				if err = parseFromReflectValue(fieldVal, sel, config); err != nil {
					return fmt.Errorf("Error parsing value to field '%s' : %s", fieldType.Name, err.Error())
				}
			}

			continue
		}
		htmlxTags, err := initializeHtmlxTags(fieldType)
		if err != nil {
			if config.parseAllFields {
				return fmt.Errorf("Error extracting tags from field '%s': %s", fieldType.Name, err.Error())
			}

			continue
		}

		htmlElement := sel.Find(htmlxTags.selector)

		if config.noEmptySelection && htmlElement.Length() == 0 {
			return fmt.Errorf("Error locating html element for field '%s'", fieldType.Name)
		}

		rawVal, err := getRawValue(fieldType, htmlElement, htmlxTags.source, config)
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
	config *Config,
) (string, error) {
	if source == "content" {
		return htmlElement.Clone().Children().Remove().End().Text(), nil
	} else if regexp.MustCompile(`^attr=[a-zA-Z-0-9]+$`).MatchString(source) {
		var exists bool
		attrName := source[5:]
		value, exists := htmlElement.Attr(attrName)
		if !exists {
			if config.noMissingAttributes {
				return "", fmt.Errorf("Error locating attribute %s for field '%s'", attrName, fieldType.Name)
			}

			return "", nil
		}

		return value, nil
	} else {
		return "", fmt.Errorf("unrecognizable source %s for field '%s'", source, fieldType.Name)
	}
}

func parseValueWithCustomParser(
	fieldVal reflect.Value,
	rawVal string,
	config *Config,
	parserName string,
) error {
	parser, ok := config.parsers[parserName]
	if !ok {
		return fmt.Errorf("parser %s is not recognizable", parserName)
	}

	val, err := parser(rawVal)
	if err != nil {
		return fmt.Errorf("parser '%s' error: %s", parserName, err.Error())
	}

	if val == nil {
		fieldVal.Set(reflect.Zero(fieldVal.Type()))
		return nil
	}

	processedVal := reflect.ValueOf(val)
	if !processedVal.IsValid() {
		return fmt.Errorf("processed value using parser %s is invalid", parserName)
	}

	if !processedVal.Type().AssignableTo(fieldVal.Type()) {
		return fmt.Errorf(
			"Incompatible type when using parser %s, want %s, get %s",
			parserName,
			fieldVal.Type().String(),
			processedVal.Type().String(),
		)
	}

	fieldVal.Set(processedVal)

	return nil
}

func parseSupportedValues(fieldVal reflect.Value, rawVal string, config *Config, htmlxTags HtmlxTags) error {
	switch fieldVal.Type() {
	case reflect.TypeOf(""):
		strVal, _ := StringParserClean(rawVal)
		fieldVal.Set(reflect.ValueOf(strVal))
	case reflect.TypeOf(0):
		intVal, err := IntParser(rawVal)
		if err != nil {
			return fmt.Errorf("Int parser error: %s", err.Error())
		}
		fieldVal.Set(reflect.ValueOf(intVal))
	case reflect.TypeOf(0.5):
		floatVal, err := FloatParser(rawVal)
		if err != nil {
			return fmt.Errorf("Float parser error: %s", err.Error())
		}
		fieldVal.Set(reflect.ValueOf(floatVal))
	case reflect.TypeOf(time.Time{}):
		dateVal, err := DateParser(config.dateFormat)(rawVal)
		if err != nil {
			return fmt.Errorf("Date parser error: %s", err.Error())
		}

		fieldVal.Set(reflect.ValueOf(dateVal))
	default:
		if fieldVal.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				ptr := reflect.New(fieldVal.Type().Elem())
				if err := parseValue(ptr.Elem(), rawVal, config, htmlxTags); err != nil {
					return err
				}

				fieldVal.Set(ptr)
				return nil
			} else {
				return parseValue(fieldVal.Elem(), rawVal, config, htmlxTags)
			}
		}

		return fmt.Errorf("Value of type %s is not supported", fieldVal.Type().String())
	}

	return nil
}

func parseValue(fieldVal reflect.Value, rawVal string, config *Config, htmlxTags HtmlxTags) error {
	if !fieldVal.IsValid() {
		return fmt.Errorf("Field doesn't represent a value")
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Field can't be set")
	}

	if htmlxTags.parser != "" {
		return parseValueWithCustomParser(fieldVal, rawVal, config, htmlxTags.parser)
	}

	if strings.TrimSpace(rawVal) == "" {
		// Skip the  field if the raw value is empty
		return nil
	}

	return parseSupportedValues(fieldVal, rawVal, config, htmlxTags)
}
