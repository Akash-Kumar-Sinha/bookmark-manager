package helpers

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func InitKeys() {
	signBytes, err := os.ReadFile("private.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to read private.pem: %v", err))
	}
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse private key: %v", err))
	}

	verifyBytes, err := os.ReadFile("public.pem")
	if err != nil {
		panic(fmt.Sprintf("Failed to read public.pem: %v", err))
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse public key: %v", err))
	}
}
