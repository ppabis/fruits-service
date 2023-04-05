package token

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestTokenCreate(t *testing.T) {
	PrivateKey, _ = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	// Creates a new token
	tokenString, err := CreateToken(123, "fruits", map[string]interface{}{"super": true})
	if err != nil {
		t.Error(err)
	}

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwt.MapClaims{})
	if err != nil {
		t.Error(err)
	}

	aud, err := token.Claims.GetAudience()
	if err != nil {
		t.Error(err)
	}

	if aud[0] != "service:fruits" {
		t.Errorf("audience is not correct %s", aud)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		t.Error(err)
	}

	if sub != "user:123" {
		t.Errorf("subject is not correct %s", sub)
	}

	claims := token.Claims.(*jwt.MapClaims)
	super := (*claims)["super"]

	if super.(bool) && super != true {
		t.Errorf("super is not correct %v", super)
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &PrivateKey.PublicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodES512.Name}))

	if err != nil {
		t.Error(err)
	}

}
