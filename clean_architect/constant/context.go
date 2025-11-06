package constant

type CtxKey string

const (
	CtxKeyRequestID CtxKey = "request_id"
	CtxKeyUserID    CtxKey = "user_id"
	CtxKeyUsername  CtxKey = "username"
	CtxKeyEmail     CtxKey = "email"
)
