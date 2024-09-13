package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CheckStatus(resp *http.Response) error {

	switch code := resp.StatusCode; {
	case code == http.StatusUnauthorized:
		return fmt.Errorf("unauthorized")
	case code == http.StatusForbidden:
		return errors.New("forbidden")
	case code == 422:
		var repr []byte
		bodyBytes, err := io.ReadAll(resp.Body)
		defer func() { _ = resp.Body.Close() }()
		if err != nil {
			return err
		}
		if strings.Contains(resp.Header.Get("Content-Type"), "json") && resp.StatusCode == 422 {
			var dest HTTPValidationError
			if err := json.Unmarshal(bodyBytes, &dest); err != nil {
				return err
			}
			repr, err = json.Marshal(*dest.Detail)
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("validation error: %s", repr)

	case code >= 500:
		return fmt.Errorf("server error: %d", code)
	case code < 300:
		return nil
	default:
		return fmt.Errorf("unrecognized status code %d", code)
	}
}
