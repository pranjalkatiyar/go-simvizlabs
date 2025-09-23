package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAppStoreJWT() (string, error) {
	// Load environment variables
	keyID := os.Getenv("APPSTORE_KEY_ID")
	issuerID := os.Getenv("APPSTORE_ISSUER_ID")
	bundleID := os.Getenv("APPSTORE_BUNDLE_ID") // optional
	rawKey := os.Getenv("APPSTORE_PRIVATE_KEY")

	// Replace escaped newlines with actual newlines
	privateKeyPEM := strings.ReplaceAll(rawKey, `\n`, "\n")

	// Decode PEM block
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("invalid PEM block")
	}

	// Parse PKCS8 private key
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("not an ECDSA key")
	}

	// Create JWT claims
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": issuerID,
		"iat": now,
		"exp": now + 300, // 5 minutes
		"aud": "appstoreconnect-v1",
	}
	if bundleID != "" {
		claims["bid"] = bundleID
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = keyID

	// Sign the token
	signedToken, err := token.SignedString(ecdsaKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}
