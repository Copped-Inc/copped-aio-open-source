package helper

import (
	"errors"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/dgrijalva/jwt-go"
	"net/http"

	"github.com/Copped-Inc/aio-types/secrets"
)

func GetClaim(key string, r *http.Request) (interface{}, error) {

	auth := cookies.Get(r, "authorization")
	if auth == nil {
		return nil, errors.New("no authorization cookie")
	}

	token, err := jwt.Parse(auth.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(secrets.JWT_Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("claims not mapclaims")
	}

	return claims[key], nil

}
