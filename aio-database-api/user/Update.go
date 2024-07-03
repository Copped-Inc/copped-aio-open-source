package user

import (
	"database-api/product"
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/modules"
	"golang.org/x/exp/slices"
)

func (d *Database) AddWebhook(w string) *Database {
	if d.Settings == nil {
		d.Settings = &Settings{}
	}

	d.Settings.Webhooks = append(d.Settings.Webhooks, w)
	return d
}

func (d *Database) DeleteWebhook(w string) *Database {
	if d.Settings != nil {
		for i, v := range d.Settings.Webhooks {
			if v == w {
				d.Settings.Webhooks = append(d.Settings.Webhooks[:i], d.Settings.Webhooks[i+1:]...)
				break
			}
		}
	}

	if d.Settings.Stores == "" && len(d.Settings.Webhooks) == 0 {
		d.Settings = nil
	}

	return d
}

func (d *Database) UpdateInstance(instance Instance) *Database {
	for i, in := range d.Instances {
		if in.ID == instance.ID {
			d.Instances[i] = instance
			d.UpdateSessionFromData()
			return d
		}
	}
	d.Instances = append(d.Instances, instance)
	d.UpdateSessionFromData()

	return d
}

func (d *Database) UpdateSessionFromData() *Database {
	s := func() Session {
		if d.Session == nil {
			return Session{}
		}
		return Session{Status: d.Session.Status}
	}()
	if s.Status == "Running" {
		for _, in := range d.Instances {
			if in.Status == "Running" {
				tm, err := strconv.Atoi(in.TaskMax)
				if err != nil {
					tm = 0
				}

				s.Instances = append(s.Instances, in)
				s.Tasks += tm
			}
		}
	}
	d.Session = &s
	return d
}

func (d *Database) AddWhitelist(p string) (*Database, error) {
	if !slices.Contains(d.Products, p) {
		d.Products = append(d.Products, p)
	}

	pr, err := product.Get(p)
	if err != nil {
		return d, err
	}

	for i, state := range pr.UserState {
		if state.ID == d.User.ID {
			pr.UserState[i].State = product.Whitelisted
			return d, pr.Save()
		}
	}

	pr.UserState = append(pr.UserState, product.UserState{
		ID:    d.User.ID,
		SKU:   pr.SKU,
		State: product.Whitelisted,
	})

	return d, pr.Save()
}

func (d *Database) RemoveWhitelist(p string) (*Database, error) {
	if !slices.Contains(d.Products, p) {
		d.Products = append(d.Products, p)
	}

	pr, err := product.Get(p)
	if err != nil {
		return d, err
	}

	for i, state := range pr.UserState {
		if state.ID == d.User.ID {
			pr.UserState[i].State = product.Blacklisted
			return d, pr.Save()
		}
	}

	pr.UserState = append(pr.UserState, product.UserState{
		ID:    d.User.ID,
		SKU:   pr.SKU,
		State: product.Blacklisted,
	})

	return d, pr.Save()
}

func (d *Database) NeedPassword(b bool) *Database {
	d.PassNeeded = b
	return d
}

func (d *Database) UpdateStore(store modules.Site, b bool) *Database {
	if d.Settings == nil {
		d.Settings = &Settings{}
	}

	var err error
	if !b {
		err = d.Settings.Stores.Remove(store)
	} else {
		err = d.Settings.Stores.Add(store)
	}
	if err != nil {
		panic(err)
	}

	if d.Settings.Stores == "" && len(d.Settings.Webhooks) == 0 {
		d.Settings = nil
	}

	return d
}

func (d *Database) DeleteInstance(instance Instance) *Database {
	for i, in := range d.Instances {
		if in.ID == instance.ID {
			d.Instances = append(d.Instances[:i], d.Instances[i+1:]...)
			return d
		}
	}
	return d
}

func (d *Database) UpdateSession(s Session) *Database {
	d.Session = &s
	return d
}

func (d *Data) UpdateBilling(b []Billing) *Data {
	d.Billing = b
	return d
}

func (d *Data) UpdateShipping(s []Shipping) *Data {
	d.Shipping = s
	return d
}

func (d *Database) SetCode(code string) *Database {
	d.User.Code = code
	d.User.CodeExpire = time.Now().Add(time.Hour * 24)
	return d
}

func (d *Database) ToDataResp(data *Data) (DataResp, error) {
	checkouts, err := d.GetCheckouts()
	if err != nil {
		return DataResp{}, err
	}

	return DataResp{
		Checkouts: checkouts,
		Instances: d.Instances,
		Settings:  d.Settings,
		Session:   d.Session,
		Billing:   data.Billing,
		Shipping:  data.Shipping,
		User:      d.User,
		Whitelist: d.Products,
	}, nil
}
