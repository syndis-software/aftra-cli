package openapi

import (
	"context"
	"net/http"
)

// NewSecurityProviderApiKey provides a SecurityProvider, which can solve
// the Auth challenge for api-calls.
func NewSecurityProviderApiKey(header_name string, token string) (*SecurityProviderApiKey, error) {
	return &SecurityProviderApiKey{
		api_key:     token,
		header_name: header_name,
	}, nil
}

// SecurityProviderApiKey sends a token as part of an
// Authorization: Bearer header along with a request.
type SecurityProviderApiKey struct {
	api_key     string
	header_name string
}

// Intercept will attach an Authorization header to the request
// and ensures that the bearer token is attached to the header.
func (s *SecurityProviderApiKey) Intercept(ctx context.Context, req *http.Request) error {
	req.Header.Set(s.header_name, s.api_key)
	return nil
}
