package models

import ct "social-network/shared/go/customtypes"

type CreatePrivateConvParams struct {
	UserA ct.Id
	UserB ct.Id
}

type CreateGroupConvParams struct {
	GroupId ct.Id
	UserIds ct.Ids
}
