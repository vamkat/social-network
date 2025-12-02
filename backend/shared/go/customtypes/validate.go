package customtypes

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidateStruct iterates over exported struct fields and validates them.
// - If a field implements Validator, its Validate() method is called.
// - If a field does not have `validate:"nullable"` tag, zero values are flagged as errors.
// - Nullable fields if empty return nil error.
// - Primitive types are excluded
// Example:
//
//	type RegisterRequest struct {
//	    Username  customtypes.Username `json:"username,omitempty" validate:"nullable"` // optional
//	    FirstName customtypes.Name     `json:"first_name,omitempty" validate:"nullable"` // optional
//	    LastName  customtypes.Name     `json:"last_name"` // required
//	    About     customtypes.About    `json:"about"`     // required
//	    Email     customtypes.Email    `json:"email,omitempty" validate:"nullable"` // optional
//	}
// func ValidateStruct(s any) error {
// 	return ValidateType(reflect.ValueOf(s), "")
// }

// func ValidateType(v reflect.Value, fieldName string) error {
// 	if v.Kind() == reflect.Pointer {
// 		if v.IsNil() {
// 			return nil
// 		}
// 		v = v.Elem()
// 	}

// 	switch v.Kind() {
// 	case reflect.Struct:
// 		return validateStruct(v)
// 	case reflect.Slice:
// 		return validateSlice(v, fieldName)
// 	default:
// 		return validateValidator(v, fieldName)
// 	}
// }

// // validateStruct checks each exported field of a struct for required rules and,
// // if applicable, runs Validator-based validation on the field value.
// func validateStruct(v reflect.Value) error {
// 	t := v.Type()
// 	var allErrors []string

// 	for i := 0; i < v.NumField(); i++ {
// 		fieldVal := v.Field(i)
// 		fieldType := t.Field(i)

// 		if !fieldVal.CanInterface() {
// 			continue
// 		}

// 		tag := fieldType.Tag.Get("validate")
// 		nullable := tag == "nullable"

// 		isPrimitive := isPrimitiveField(fieldVal)
// 		_, zeroOk := allowedZeroVal[fieldVal.Type().Name()]

// 		// Required-field check
// 		if !nullable && !isPrimitive && !zeroOk && isZeroValue(fieldVal) {
// 			allErrors = append(allErrors,
// 				fmt.Sprintf("%s: required field missing", fieldType.Name))
// 			continue
// 		}

// 		err := validateValidator(fieldVal, fieldType.Name)

// 		if err != nil {
// 			allErrors = append(allErrors, err.Error())
// 		}
// 	}

// 	if len(allErrors) > 0 {
// 		return fmt.Errorf("validation errors: %v", allErrors)
// 	}
// 	return nil
// }

// func validateSlice(v reflect.Value, fieldName string) error {

// 	elemType := v.Type().Elem()
// 	isPrimitive := elemType.PkgPath() == ""

// 	// Primitive slices require no element validation
// 	if isPrimitive {
// 		return nil
// 	}

// 	var errs []string

// 	for i := 0; i < v.Len(); i++ {
// 		elemVal := v.Index(i)
// 		elemAny := elemVal.Interface()

// 		validator, ok := elemAny.(Validator)
// 		if !ok {
// 			errs = append(errs,
// 				fmt.Sprintf("%s[%d]: element does not implement Validator", fieldName, i))
// 			continue
// 		}

// 		if err := validator.Validate(); err != nil {
// 			errs = append(errs,
// 				fmt.Sprintf("%s[%d]: %v", fieldName, i, err))
// 		}
// 	}

// 	if len(errs) > 0 {
// 		return fmt.Errorf("validation errors: %v", errs)
// 	}
// 	return nil
// }

// // validateValidator runs validation on a value only if it implements
// // the Validator interface; primitives and non-validator types are ignored.
// func validateValidator(v reflect.Value, fieldName string) error {
// 	val := v.Interface()

// 	validator, ok := val.(Validator)
// 	if !ok {
// 		return nil // primitive or non-validator → no extra validation
// 	}

// 	if err := validator.Validate(); err != nil {
// 		return fmt.Errorf("%s: %v", fieldName, err)
// 	}
// 	return nil
// }

// func isPrimitiveField(v reflect.Value) bool {
// 	if v.Kind() == reflect.Slice {
// 		return v.Type().Elem().PkgPath() == ""
// 	}
// 	return v.Type().PkgPath() == ""
// }

func ValidateStruct(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	var allErrors []string

	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		val := fieldVal.Interface()
		validator, ok := val.(Validator)

		validateTag := fieldType.Tag.Get("validate")
		// Check for exact "nullable" match (not "elements=nullable")
		nullable := validateTag == "nullable" ||
			(strings.Contains(validateTag, "nullable") && !strings.HasPrefix(validateTag, "elements="))
		elementsNullable := strings.Contains(validateTag, "elements=nullable")
		isPrimitive := false

		if fieldVal.Kind() != reflect.Slice {
			// Normal fields: primitive if they have no package
			isPrimitive = fieldVal.Type().PkgPath() == ""
		} else {
			// Slice fields: primitive only if the element type is primitive
			elemType := fieldVal.Type().Elem()
			isPrimitive = elemType.PkgPath() == ""
		}
		_, zeroOk := allowedZeroVal[fieldVal.Type().Name()]

		if !nullable && !isPrimitive && !zeroOk {
			if isZeroValue(fieldVal) {
				allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
				continue
			}
		}

		// --- SLICE VALIDATION ---
		if fieldVal.Kind() == reflect.Slice {

			// Treat empty slice as nil for required check
			if fieldVal.IsNil() || fieldVal.Len() == 0 {
				if !nullable {
					allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
				}
				continue
			}

			// Skip primitive slices for element validation
			if isPrimitive {
				continue
			}

			// Validate each element (customtypes / Validator)
			for j := 0; j < fieldVal.Len(); j++ {
				elem := fieldVal.Index(j).Interface()
				elemVal := reflect.ValueOf(elem)

				if vElem, ok := elem.(Validator); ok {
					// Check if element is zero/empty when elements are required
					if !elementsNullable && isZeroValue(elemVal) {
						allErrors = append(allErrors,
							fmt.Sprintf("%s[%d]: required element missing", fieldType.Name, j))
						continue
					}

					if err := vElem.Validate(); err != nil {
						allErrors = append(allErrors,
							fmt.Sprintf("%s[%d]: %v", fieldType.Name, j, err))
					}
				} else {
					// Non-Validator element in a non-primitive slice
					allErrors = append(allErrors,
						fmt.Sprintf("%s[%d]: element does not implement Validator", fieldType.Name, j))
				}
			}
			continue
		}

		if ok {
			if err := validator.Validate(); err != nil {
				allErrors = append(allErrors, fmt.Sprintf("%s: %v", fieldType.Name, err))
			}
		}
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("validation errors: %v", allErrors)
	}
	return nil
}

// isZeroValue returns true if the reflect.Value is its type's zero value
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
