package htmlx

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Parse the content from HTML string to struct
func ParseFromString(
	s any,
	content string,
) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return err
	}
	return ParseFromSelection(s, doc.Selection)
}

// Parse the content from goquery Selection to struct
func ParseFromDocument(
	s any,
	doc *goquery.Document,
) error {
	return ParseFromSelection(s, doc.Selection)
}

// Parse the content from HTML string to struct with context
func ParseFromStringWithContext(
	ctx context.Context,
	s any,
	content string,
) error {
	errChan := make(chan error)

	go func() {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			errChan <- err
		}

		errChan <- ParseFromSelection(s, doc.Selection)
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
) error {
	errChan := make(chan error)

	go func() {
		errChan <- ParseFromSelection(s, doc.Selection)
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
) error {
	errChan := make(chan error)

	go func() {
		errChan <- ParseFromSelection(s, sel)
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
) error {
	v := reflect.ValueOf(s)

	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("s must be a pointer to a struct")
	}

	v = v.Elem()

	if err := parseFromReflectValue(v, sel); err != nil {
		return err
	}

	return nil
}

func parseFromReflectValue(
	v reflect.Value,
	sel *goquery.Selection,
) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("value is not a struct")
	}
	t := v.Type()

	for i := range t.NumField() {
		var err error

		fieldType := t.Field(i)
		fieldVal := v.Field(i)

		if fieldVal.Kind() == reflect.Struct {
			if err = parseFromReflectValue(fieldVal, sel); err != nil {
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

		if err = parseValue(&fieldVal, rawVal); err != nil {
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

func parseValue(fieldVal *reflect.Value, rawVal string) error {
	if !fieldVal.IsValid() {
		return fmt.Errorf("Field doesn't represent a value")
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Field can't be set")
	}

	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(strings.TrimSpace(rawVal))
	case reflect.Int:
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
	case reflect.Float64:
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
	}

	return nil
}
