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
	if err := ValidateStruct(k); err != nil {
		return "", err
	}
	return fmt.Sprintf("img_%s:%d", k.Variant, k.Id), nil
}

type IsGroupMemberKey struct {
	GroupId Id
	UserId  Id
}

func (im IsGroupMemberKey) String() (string, error) {
	if err := ValidateStruct(im); err != nil {
		return "", err
	}
	return fmt.Sprintf("is_group:%d.member:%d", im.GroupId.Int64(), im.UserId.Int64()), nil
}
