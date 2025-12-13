package contextkeys

import (
	ct "social-network/shared/go/customtypes"
	"social-network/shared/go/gorpc"
)

var commonKeys = []gorpc.StringableKey{ct.UserId, ct.ReqID, ct.TraceId}

func CommonKeys(extraKeys ...gorpc.StringableKey) []gorpc.StringableKey {
	return append(commonKeys, extraKeys...)
}
