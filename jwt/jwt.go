package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Claims is a map of string->something containing the meta infos associated with a token
type Claims map[string]interface{}

// CreateToken takes some claims and a private key (either rsa or ec) and returns a signed json web token
func CreateToken(claims Claims, key interface{}) (string, error) {
	switch k := key.(type) {
	case *rsa.PrivateKey:
		{
			token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims(claims))
			return token.SignedString(k)
		}
	case *ecdsa.PrivateKey:
		{
			token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims(claims))
			return token.SignedString(k)
		}
	}
	return "", errors.New("invalid private key")
}

// ValidateToken checks the signature of the token with a given public key and returns the associated claims
func ValidateToken(tokenString string, key interface{}) (Claims, error) {
	var (
		token *jwt.Token
		err   error
	)
	switch k := key.(type) {
	case *rsa.PublicKey:
		{
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return k, nil
			})
		}
	case *ecdsa.PublicKey:
		{
			token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return k, nil
			})
		}
	}
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return Claims(claims), nil
	}
	return nil, errors.New("invalid token")
}

// LoadPublicKey loads a PEM encoded public key (either rsa or ec)
func LoadPublicKey(keyFile string) (interface{}, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	rsaKey, err := jwt.ParseRSAPublicKeyFromPEM(bs)
	if err != nil {
		ecKey, err := jwt.ParseECPublicKeyFromPEM(bs)
		if err != nil {
			return nil, errors.New("unknown public key type")
		}
		return ecKey, nil
	}
	return rsaKey, nil
}

// LoadPrivateKey loads a PEM encoded private key (either rsa or ec)
func LoadPrivateKey(keyFile string) (interface{}, error) {
	bs, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM(bs)
	if err != nil {
		ecKey, err := jwt.ParseECPrivateKeyFromPEM(bs)
		if err != nil {
			return nil, errors.New("unknown public key type")
		}
		return ecKey, nil
	}
	return rsaKey, nil
}

// GetTokenFromRequest takes the first Authorization header and extracts the bearer json web token
func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("no valid authorization header")
	}
	return parts[1], nil
}

// GetClaimsFromRequest extracts the
func GetClaimsFromRequest(r *http.Request, key interface{}) (Claims, error) {
	token, err := GetTokenFromRequest(r)
	if err != nil {
		return nil, err
	}
	return ValidateToken(token, key)
}
