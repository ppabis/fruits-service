package token

import (
	"fmt"
	"monolith/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Creates a token for a user.
// `user` is the user id (subject = user:123),
// `service` is the name of the service (audience = service:fruits),
// `params` are extras, like is user super,
// It sets 30 seconds as the standard expiration time.
func CreateToken(user int, service string, params map[string]interface{}) (string, error) {

	// Fill claims with params
	claims := jwt.MapClaims{}

	for k, v := range params {
		claims[k] = v
	}

	// Overwrite arguments that are standard
	claims["sub"] = fmt.Sprintf("user:%d", user)
	claims["aud"] = fmt.Sprintf("service:%s", service)
	claims["exp"] = time.Now().Add(config.TokenLifetime).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	return token.SignedString(config.PrivateKey)
}
