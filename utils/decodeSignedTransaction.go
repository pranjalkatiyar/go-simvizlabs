package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func DecodeSignedTransactionInfo(jws string) (map[string]interface{}, error) {
	parts := strings.Split(jws, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWS format")
	}

	payload := parts[1]

	// Base64URL decode
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	// Unmarshal JSON
	var result map[string]interface{}
	if err := json.Unmarshal(decoded, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return result, nil
}
