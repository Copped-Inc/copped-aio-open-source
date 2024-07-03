package user

import (
	"bytes"
	"database-api/database"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/subscriptions"
)

func FromRequest(r *http.Request) (*Database, error) {
	id, err := helper.GetClaim("id", r)
	if err != nil || id == "" {
		return nil, err
	}

	doc, err := database.GetDatabase().Collection("data").Doc(fmt.Sprintf("%s", id)).Get(database.GetContext())
	if err != nil {
		return nil, err
	}

	var data Database
	err = doc.DataTo(&data)
	return &data, err
}

func FromId(id string) (Database, error) {
	doc, err := database.GetDatabase().Collection("data").Doc(id).Get(database.GetContext())
	if err != nil {
		return Database{}, err
	}

	var data Database
	err = doc.DataTo(&data)
	return data, err
}

func DataFromRequest(r *http.Request) (*Data, *Database, error) {

	data, err := FromRequest(r)
	if err != nil {
		return nil, nil, err
	}

	if !bytes.Equal(data.Password, helper.CreateHash(r.Header.Get("password"))) {
		return nil, nil, errors.New("password is incorrect")
	}

	j, err := helper.Decrypt(data.Data, r.Header.Get("password"))
	if err != nil {
		return nil, nil, err
	}

	var d Data
	err = json.Unmarshal(j, &d)
	return &d, data, err

}

func DataFromWebsocket(id, p string) (*Data, *Database, error) {

	doc, err := database.GetDatabase().Collection("data").Doc(id).Get(database.GetContext())
	if err != nil {
		return nil, nil, err
	}

	var data Database
	err = doc.DataTo(&data)
	if err != nil {
		return nil, nil, err
	}

	if !bytes.Equal(data.Password, helper.CreateHash(p)) {
		return nil, nil, errors.New("password is incorrect")
	}

	j, err := helper.Decrypt(data.Data, p)
	if err != nil {
		return nil, nil, err
	}

	var d Data
	err = json.Unmarshal(j, &d)
	return &d, &data, err

}

func GetAll() ([]Database, error) {

	var d []Database
	docs, err := database.GetDatabase().Collection("data").Documents(database.GetContext()).GetAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		var u Database
		if err = doc.DataTo(&u); err != nil {
			return nil, err
		}

		d = append(d, u)
	}

	return d, err
}

func GetWithPw(f func(http.ResponseWriter, *http.Request, *Data, *Database)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d, db, err := DataFromRequest(r)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusBadRequest)
			return
		}

		f(w, r, d, db)
	}
}

func Get(f func(http.ResponseWriter, *http.Request, *Database)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := FromRequest(r)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusBadRequest)
			return
		}

		f(w, r, db)
	}
}

func Verify(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := FromRequest(r)
		if err != nil {
			VerifyAdmin(f)(w, r)
			return
		}

		f(w, r)
	}
}

func VerifyAdmin(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !helper.IsMaster(r.Header.Get("Password")) {
			if db, err := FromRequest(r); err != nil {
				console.ErrorRequest(w, r, err, http.StatusUnauthorized)
				return
			} else if db.User.Subscription.Plan != subscriptions.Developer {
				console.ErrorRequest(w, r, err, http.StatusUnauthorized)
				return
			}
		}

		f(w, r)
	}
}
