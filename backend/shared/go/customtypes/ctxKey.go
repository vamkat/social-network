package customtypes

// type alias:
// context key type ALIAS in order to help enforcing a single source of truth for key namings
type CtxKey string

func (c CtxKey) String() string {
	return string(c)
}

// Holds the keys to values on request context.
// warning! keys that are meant to be propagated through grpc services have strict requirements! They must be ascii, lowercase, and only allowed symbols: "-_."
const (
	ClaimsKey        CtxKey = "jwt-claims"
	ReqActionDetails CtxKey = "action-details"
	ReqTimestamp     CtxKey = "timestamp"
	ReqID            CtxKey = "request-id"
	UserId           CtxKey = "user-id"
	TraceId          CtxKey = "trace-id"
)
