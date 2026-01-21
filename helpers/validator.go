package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidatePayload unmarshals and validates payload
func ValidatePayload(cBody []byte, req interface{}) string {
	// Step 1: unmarshal JSON
	err := json.Unmarshal(cBody, req)
	if err != nil {
		// jika type mismatch
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			fieldName := ute.Field
			return fmt.Sprintf("%s must be %s", fieldName, ute.Type.Name())
		}
		return err.Error()
	}

	// Step 2: validate struct
	err = validate.Struct(req)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			firstErr := validationErrors[0]

			// ambil json key
			rt := reflect.TypeOf(req).Elem()
			field, _ := rt.FieldByName(firstErr.StructField())
			jsonKey := field.Tag.Get("json")
			if jsonKey == "" {
				jsonKey = firstErr.StructField()
			}

			// custom message based on tag
			var msg string
			switch firstErr.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", jsonKey)
			case "min":
				msg = fmt.Sprintf("%s minimum length %s", jsonKey, firstErr.Param())
			case "gt":
				msg = fmt.Sprintf("%s must be greater than %s", jsonKey, firstErr.Param())
			default:
				msg = firstErr.Error()
			}
			return msg
		}
		return err.Error()
	}

	return ""
}
