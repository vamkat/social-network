package ct

import (
	"fmt"
)

type BasicUserInfoKey struct {
	Id Id
}

func (k BasicUserInfoKey) String() (string, error) {
	if err := k.Id.Validate(); err != nil {
		return "", err
	}
	return fmt.Sprintf("basic_user_info:%d", k.Id), nil
}

type ImageKey struct {
	Variant FileVariant
	Id      Id
}

func (k ImageKey) String() (string, error) {
	if err := k.Variant.Validate(); err != nil {
		return "", err
	}
	if err := k.Id.Validate(); err != nil {
		return "", err
	}
	return fmt.Sprintf("img_%s:%d", k.Variant, k.Id), nil
}
