package constant

type ContextKey string

const LoggerKey = ContextKey("logger")
const RequestIDKey = ContextKey("request_id")
