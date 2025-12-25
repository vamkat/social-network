package ct

import (
	"fmt"
	"reflect"
)

// ValidateStruct iterates over exported struct fields and validates them.
//   - If a field implements Validator, its Validate() method is called.
//   - If a field does not have `validate:"nullable"` tag, zero values are flagged as errors.
//   - Nullable fields if empty return nil error.
//   - All primitives are excluded except slices containing custom types.
//   - If a field is a slice of a custom type null values in that slice are considered invalid.
//
// Example:
//
//	type RegisterRequest struct {
//		    Username  ct.Username 		`json:"username,omitempty" validate:"nullable"` // optional
//		    FirstName ct.Name     		`json:"first_name,omitempty" validate:"nullable"` // optional
//		    LastName  ct.Name     		`json:"last_name"` // required
//		    About     ct.About    		`json:"about"`     // required
//		    Email     ct.Email    		`json:"email,omitempty" validate:"nullable"` // optional
//	}
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
		validator, implementsValidator := val.(Validator)
		if !implementsValidator {
			continue
		}

		validateTag := fieldType.Tag.Get("validate")
		nullable := validateTag == "nullable"
		_, zeroOk := alwaysAllowZero[fieldVal.Type().Name()]

		// Handle custom types with zero (null) value
		if isZeroValue(fieldVal) {
			// Check for null fields that are not allowed to be null
			if !nullable && !zeroOk {
				allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
				continue
			}

			// Skip validation for nullable fields that are empty
			if nullable {
				continue
			}
		}

		if implementsValidator {
			if err := validator.Validate(); err != nil {
				allErrors = append(allErrors, fmt.Sprintf("%s: %v", fieldType.Name, err))
			}
			continue
		}

		// --- SLICE VALIDATION ---
		// if fieldVal.Kind() == reflect.Slice {

		// 	// Treat empty slice as nil for required check
		// 	if fieldVal.IsNil() || fieldVal.Len() == 0 {
		// 		if !nullable {
		// 			allErrors = append(allErrors, fmt.Sprintf("%s: required field missing", fieldType.Name))
		// 		}
		// 		continue
		// 	}

		// 	// Skip primitive slices for element validation
		// 	if isPrimitive {
		// 		continue
		// 	}

		// 	// Validate each element (customtypes / Validator)
		// 	for j := 0; j < fieldVal.Len(); j++ {
		// 		elem := fieldVal.Index(j).Interface()
		// 		elemVal := reflect.ValueOf(elem)

		// 		if vElem, ok := elem.(Validator); ok {
		// 			// Check if element is zero/empty when elements are required
		// 			if isZeroValue(elemVal) {
		// 				allErrors = append(allErrors,
		// 					fmt.Sprintf("%s[%d]: required element missing", fieldType.Name, j))
		// 				continue
		// 			}

		// 			if err := vElem.Validate(); err != nil {
		// 				allErrors = append(allErrors,
		// 					fmt.Sprintf("%s[%d]: %v", fieldType.Name, j, err))
		// 			}
		// 		} else {
		// 			// Non-Validator element in a non-primitive slice
		// 			allErrors = append(allErrors,
		// 				fmt.Sprintf("%s[%d]: element does not implement Validator", fieldType.Name, j))
		// 		}
		// 	}
		// 	continue
		// }
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("validation errors: %v", allErrors)
	}
	return nil
}

// isZeroValue returns true if the reflect.Value is its type's zero value
func isZeroValue(v reflect.Value) bool {
	// If it's a slice, treat len == 0 as zero value
	if v.Kind() == reflect.Slice {
		return v.IsNil() || v.Len() == 0
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
