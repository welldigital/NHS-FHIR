package middleware

import (
	"net/http"
)

func ValidateToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO: middleware to validate the token
	// verify against public key?
	// make sure not expired

	// validate token
	// token, err := jwt.ParseWithClaims(tokenSigned, claims, func(t *jwt.Token) (interface{}, error) {
	// 	t.Header["kid"] = "test-1"
	// 	return secretKey, nil
	// })

}
