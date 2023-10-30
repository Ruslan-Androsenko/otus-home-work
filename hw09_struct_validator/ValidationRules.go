package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	typeMin    = "min"
	typeMax    = "max"
	typeLen    = "len"
	typeIn     = "in"
	typeRegexp = "regexp"
)

type ValidationRules struct {
	Min           int
	Max           int
	Len           int
	Regexp        string
	InRange       []string
	Types         []string
	ErrValidators []error
}

// Сформировать ошибки для правил валидации.
func (rules ValidationRules) makeRulesErrors(fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.ErrValidators))

	for _, validationErr := range rules.ErrValidators {
		resErrors = append(resErrors, ValidationError{
			Field: fieldName,
			Err:   validationErr,
		})
	}

	return resErrors
}

// Валидация целочисленных значений.
func (rules ValidationRules) validateInt(value int, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.Types))

	for _, validationType := range rules.Types {
		var errMessage error

		switch validationType {
		case typeMin:
			hasValid := value >= rules.Min
			if !hasValid {
				errMessage = ErrValueIsLessThanSpecified
			}

		case typeMax:
			hasValid := value <= rules.Max
			if !hasValid {
				errMessage = ErrValueIsGreaterThanSpecified
			}

		case typeIn:
			hasValid, err := sliceContains(value, reflect.Int, rules.InRange)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrValueOutsideSpecifiedRange
			}
		}

		if errMessage != nil {
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return resErrors
}

// Валидация строковых значений.
func (rules ValidationRules) validateString(value string, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(rules.Types))

	for _, validationType := range rules.Types {
		var errMessage error

		switch validationType {
		case typeLen:
			hasValid := len(value) == rules.Len
			if !hasValid {
				errMessage = ErrLengthDoesNotMatchSpecified
			}

		case typeIn:
			hasValid, err := sliceContains(value, reflect.String, rules.InRange)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrValueOutsideSpecifiedRange
			}

		case typeRegexp:
			hasValid, err := regexp.MatchString(rules.Regexp, value)
			if err != nil {
				resErrors = append(resErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}

			if !hasValid {
				errMessage = ErrValueDoesNotMatchRegExp
			}
		}

		if errMessage != nil {
			resErrors = append(resErrors, ValidationError{
				Field: fieldName,
				Err:   errMessage,
			})
		}
	}

	return resErrors
}

// Валидация целочисленных слайсов.
func (rules ValidationRules) validateSliceInt(items []int, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errsValidate := rules.validateInt(itemValue, fieldName)
		if len(errsValidate) > 0 {
			resErrors = append(resErrors, errsValidate...)
		}
	}

	return resErrors
}

// Валидация строковых слайсов.
func (rules ValidationRules) validateSliceString(items []string, fieldName string) []ValidationError {
	resErrors := make([]ValidationError, 0, len(items))

	for _, itemValue := range items {
		errsValidate := rules.validateString(itemValue, fieldName)
		if len(errsValidate) > 0 {
			resErrors = append(resErrors, errsValidate...)
		}
	}

	return resErrors
}

// Сформировать правила для валидации.
func makeValidators(tagValidate string) ValidationRules {
	var rules ValidationRules
	validators := strings.Split(tagValidate, "|")

	for _, validator := range validators {
		validationType := strings.Split(validator, ":")[0]

		switch validationType {
		case typeMin:
			value, err := castNumberValidator(validator, typeMin)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Min = value
			rules.Types = append(rules.Types, typeMin)

		case typeMax:
			value, err := castNumberValidator(validator, typeMax)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Max = value
			rules.Types = append(rules.Types, typeMax)

		case typeLen:
			value, err := castNumberValidator(validator, typeLen)
			if err != nil {
				rules.ErrValidators = append(rules.ErrValidators, err)
				continue
			}

			rules.Len = value
			rules.Types = append(rules.Types, typeLen)

		case typeIn:
			rules.InRange = strings.Split(getValidatorValue(validator, typeIn), ",")
			rules.Types = append(rules.Types, typeIn)

		case typeRegexp:
			rules.Regexp = getValidatorValue(validator, typeRegexp)
			rules.Types = append(rules.Types, typeRegexp)
		}
	}

	return rules
}

// Привести значение правила валидатора к числу.
func castNumberValidator(validator, validatorType string) (int, error) {
	var resErr error
	value, err := strconv.Atoi(getValidatorValue(validator, validatorType))
	if err != nil {
		resErr = ErrCanNotCastValueToNumber
	}

	return value, resErr
}

// Получить значение для валидатора.
func getValidatorValue(validator, replacement string) string {
	return strings.Replace(validator, replacement+":", "", 1)
}

// Содержится ли значение в слайсе.
func sliceContains(value interface{}, kind reflect.Kind, inRanges []string) (bool, error) {
	for _, itemRange := range inRanges {
		switch kind { //nolint:exhaustive
		case reflect.Int:
			itemRangeValue, err := strconv.Atoi(itemRange)
			if err != nil {
				return false, err
			}

			if value == itemRangeValue {
				return true, nil
			}

		case reflect.String:
			if value == itemRange {
				return true, nil
			}
		}
	}

	return false, nil
}
