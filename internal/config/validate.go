package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validateConfig(cfg any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(cfg)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	ok := errors.As(err, &validationErrors)
	if !ok {
		return err
	}

	errs := make([]error, 0)
	for _, validationErr := range validationErrors {
		errs = append(errs, buildValidationError(cfg, validationErr))
	}

	return fmt.Errorf("param errors:\n%w", errors.Join(errs...))
}

func buildValidationError(cfg any, err validator.FieldError) error {
	fieldNames := strings.Split(err.StructNamespace(), ".")[1:]
	typeOf := reflect.TypeOf(cfg).Elem()
	for index, fieldName := range fieldNames {
		structField, ok := typeOf.FieldByName(fieldName)
		if !ok {
			return &validator.InvalidValidationError{Type: err.Type()}
		}
		typeOf = structField.Type
		fieldNames[index] = getFieldName(structField, fieldName)
	}

	errMessageBuilder := strings.Builder{}
	errMessageBuilder.WriteString(fmt.Sprintf(
		"param: '%v'\terror: field validation failed on the '%s",
		strings.Join(fieldNames, "."),
		err.ActualTag(),
	))
	if err.Param() == "" {
		errMessageBuilder.WriteString("' tag")
	} else {
		errMessageBuilder.WriteString(fmt.Sprintf("=%s' tag", err.Param()))
	}

	return errors.New(errMessageBuilder.String())
}

func getFieldName(structField reflect.StructField, fieldName string) string {
	result := structField.Tag.Get("mapstructure")
	if result == "" {
		result = strings.ToLower(fieldName)
	}

	return result
}
