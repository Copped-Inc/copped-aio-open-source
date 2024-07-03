package user

import (
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (u User) Jwt() (string, error) {

	Claims := jwt.MapClaims{}

	Claims["id"] = u.ID
	Claims["created"] = time.Now().Unix()

	return jwt.NewWithClaims(jwt.SigningMethodHS512, Claims).SignedString([]byte(secrets.JWT_Secret))

}
