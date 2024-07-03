package helper

import "github.com/Copped-Inc/aio-types/secrets"

func IsMaster(password string) bool {

	return password == secrets.API_Admin_PW

}
