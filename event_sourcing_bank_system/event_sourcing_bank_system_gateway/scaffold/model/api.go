package model

type Api struct {
	Name         string
	Accept       string
	Method       string
	Path         string
	Permission   string
	RequestType  string
	ResponseType string
}

var APIPublic = map[string]bool{
	// example
	// "/api/v1/auth/user/login":                              true,
}

var PathPublic = map[string]bool{
	// example
	// "/api/v1/auth/user":                     true,
}

var ServicePath = map[string]string{
	"/api/v1/user-service": "UserService",
}
