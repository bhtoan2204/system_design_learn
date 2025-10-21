package docs

import _ "embed"

// OpenAPI holds the swagger specification for the HTTP layer.
//
//go:embed openapi.yaml
var OpenAPI []byte
