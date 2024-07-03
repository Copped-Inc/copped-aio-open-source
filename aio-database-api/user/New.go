package user

import (
	"bytes"
	"database-api/database"
	"database-api/mail"
	"encoding/json"
	"errors"

	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"

	"net/http"

	"github.com/Copped-Inc/aio-types/user"

	"github.com/Copped-Inc/aio-types/subscriptions"
)

func New(resp helper.DiscordMeResp, plan subscriptions.Plan, instanceLimit int) (*Database, error) {

	_, err := database.GetDatabase().Collection("data").Doc(resp.Id).Get(database.GetContext())
	if err == nil {
		return nil, errors.New("user already exists")
	}

	data := Database{
		User: User{
			Name:    resp.Username,
			Email:   resp.Email,
			ID:      resp.Id,
			Picture: resp.Avatar,
			Subscription: user.Subscription{
				Plan: plan,
				State: func() user.State {
					if plan != subscriptions.Developer {
						return user.Pending
					}
					return user.Active
				}(),
			},
			InstanceLimit: instanceLimit,
		},
	}

	if _, err = database.GetDatabase().Collection("data").Doc(data.User.ID).Set(database.GetContext(), data); err != nil {
		return nil, err
	}

	mail.New().
		SetTitle("Welcome to Copped AIO").
		SetSubtitle("Your account has been successfully created with the \""+plan.GetData().Name+"\" plan. You can now access the Dashboard.").
		SetButton("Dashboard", "https://aio.copped-inc.com").
		Send(data.User.Email)

	return &data, err
}

func (d *Database) Update() error {
	go AddUpdate()
	_, err := database.GetDatabase().Collection("data").Doc(d.User.ID).Set(database.GetContext(), d)
	return err
}

func AddUpdate() {
	ref := realtimedb.GetDatabase().NewRef("/userstats/global/user/updates")
	var updates int
	err := ref.Get(realtimedb.GetContext(), &updates)
	if err != nil {
		return
	}

	err = ref.Set(realtimedb.GetContext(), updates+1)
	if err != nil {
		return
	}
}

func (d *Data) UpdateDataWithPassword(p string, database *Database) error {

	if !bytes.Equal(database.Password, helper.CreateHash(p)) {
		return errors.New("password is incorrect")
	}

	j, err := json.Marshal(d)
	if err != nil {
		return err
	}

	encrypted, err := helper.Encrypt(j, p)
	if err != nil {
		return err
	}

	database.Data = encrypted
	err = database.Update()
	return err

}

func (d *Data) UpdateData(r *http.Request, database *Database) error {

	return d.UpdateDataWithPassword(r.Header.Get("password"), database)

}
