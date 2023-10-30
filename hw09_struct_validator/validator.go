package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrValueIsLessThanSpecified    = errors.New("value is less than specified")
	ErrValueIsGreaterThanSpecified = errors.New("value is greater than specified")
	ErrLengthDoesNotMatchSpecified = errors.New("length does not match specified")
	ErrValueOutsideSpecifiedRange  = errors.New("value outside the specified range")
	ErrValueDoesNotMatchRegExp     = errors.New("value does not match regular expression")
	ErrCanNotCastValueToNumber     = errors.New("can't cast value to number")
	ErrValueIsNotStructure         = errors.New("value is not a structure")
	ErrDataTypeIsNotSupported      = errors.New("data type is not supported")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}

	for _, validationErr := range v {
		errMessage := fmt.Sprintf("%s: %v; ", validationErr.Field, validationErr.Err)
		builder.WriteString(errMessage)
	}

	return builder.String()
}

func Validate(v interface{}) error {
	itemType := reflect.TypeOf(v)
	if itemType.Kind() != reflect.Struct {
		return ErrValueIsNotStructure
	}

	itemValue := reflect.ValueOf(v)
	numFields := itemType.NumField()
	resErrors := make([]ValidationError, 0, numFields)

	for i := 0; i < numFields; i++ {
		fieldOfType := itemType.Field(i)
		fieldValue := itemValue.Field(i)
		fieldName := fieldOfType.Name

		// Проверяем является ли поле публичным
		if !fieldOfType.IsExported() {
			continue
		}

		// Получаем теги для текщего поля
		fieldTag := fieldOfType.Tag
		tagValidate, ok := fieldTag.Lookup("validate")
		if !ok || len(fieldTag) == 0 {
			continue
		}

		// Формируем массив правил валидации
		rules := makeValidators(tagValidate)
		if len(rules.ErrValidators) > 0 {
			resErrors = append(resErrors, rules.makeRulesErrors(fieldName)...)
		}

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Int:
			errsValidate := rules.validateInt(int(fieldValue.Int()), fieldName)
			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		case reflect.String:
			errsValidate := rules.validateString(fieldValue.String(), fieldName)
			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		case reflect.Slice:
			var errsValidate []ValidationError

			switch items := fieldValue.Interface().(type) {
			case []int:
				errsValidate = rules.validateSliceInt(items, fieldName)

			case []string:
				errsValidate = rules.validateSliceString(items, fieldName)
			}

			if len(errsValidate) > 0 {
				resErrors = append(resErrors, errsValidate...)
			}

		default:
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   ErrDataTypeIsNotSupported,
			})
		}
	}

	if len(resErrors) > 0 {
		return ValidationErrors(resErrors)
	}

	return nil
}
