package configutil

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrBadArgument     = errors.New("argument must be pointer to struct")
	ErrUnsettableField = errors.New("field cannot be set (must be exported)")
	ErrMissingEnv      = errors.New("required env var missing")
	ErrBadConversion   = errors.New("env var conversion failed")
	ErrNoTaggedFields  = errors.New("no env-tagged fields found")
)

func LoadConfigs(localConfig any) error {
	fmt.Println("before:", localConfig)
	reflectVal := reflect.ValueOf(localConfig)
	if reflectVal.Kind() != reflect.Ptr || reflectVal.Elem().Kind() != reflect.Struct {
		return ErrBadArgument
	}

	strctVal := reflectVal.Elem()
	strctType := strctVal.Type()

	for i := 0; i < strctVal.NumField(); i++ {
		valField := strctVal.Field(i)
		typeField := strctType.Field(i)

		tagVal := typeField.Tag.Get("env")
		if tagVal == "" {
			continue
		}

		if !valField.CanSet() {
			return fmt.Errorf("%w: %s", ErrUnsettableField, typeField.Name)
		}

		envVal, ok := os.LookupEnv(tagVal)
		if !ok {
			continue
		}

		switch valField.Kind() {
		case reflect.Int:
			v, err := strconv.ParseInt(envVal, 10, 64)
			if err != nil {
				return fmt.Errorf("%w (%s): %v", ErrBadConversion, tagVal, err)
			}
			valField.SetInt(v)

		case reflect.String:
			valField.SetString(envVal)

		case reflect.Float64:
			v, err := strconv.ParseFloat(envVal, 64)
			if err != nil {
				return fmt.Errorf("%w (%s): %v", ErrBadConversion, tagVal, err)
			}
			valField.SetFloat(v)

		default:
			return fmt.Errorf("unsupported kind %s on field %s", valField.Kind(), typeField.Name)
		}

	}
	fmt.Println("after:", localConfig)
	return nil
}
