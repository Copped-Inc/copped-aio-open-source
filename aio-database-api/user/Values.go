package user

import (
	"encoding/json"
	"time"

	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/user"
)

type (
	Database struct {
		PassNeeded     bool       `json:"pass_needed" firestore:"pass_needed"`
		User           User       `json:"user" firestore:"user"`
		Data           []byte     `json:"data" firestore:"data"`
		Password       []byte     `json:"password" firestore:"password"`
		Secrets        []byte     `json:"secrets,omitempty" firestore:"secrets,omitempty"`
		CheckoutAmount int        `json:"checkout_amount" firestore:"checkout_amount"`
		Instances      []Instance `json:"instances,omitempty" firestore:"instances,omitempty"`
		Settings       *Settings  `json:"settings,omitempty" firestore:"settings,omitempty"`
		Session        *Session   `json:"session,omitempty" firestore:"session,omitempty"`
		IPs            []string   `json:"ips,omitempty" firestore:"ips,omitempty"`
		Products       []string   `json:"products,omitempty" firestore:"products,omitempty"`
	}

	User user.User

	Secrets struct {
		Oauth2 Oauth2
	}

	Data struct {
		Billing  []Billing
		Shipping []Shipping
	}

	Oauth2 struct {
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		Expiry       time.Time `json:"expiry"`
	}

	Product struct {
		User  string       `json:"id,omitempty" firestore:"user,omitempty"`
		Date  time.Time    `json:"date" firestore:"date"`
		Name  string       `json:"name" firestore:"name"`
		Link  string       `json:"link" firestore:"link"`
		Image string       `json:"image" firestore:"image"`
		Store modules.Site `json:"store" firestore:"store"`
		Size  string       `json:"size" firestore:"size"`
		Price float64      `json:"price" firestore:"price,omitempty"`
	}

	Instance struct {
		Price    float64 `json:"price" firestore:"price"`
		Provider string  `json:"provider" firestore:"provider"`
		ID       string  `json:"id" firestore:"id"`
		Status   string  `json:"status" firestore:"status"`
		TaskMax  string  `json:"task_max" firestore:"task_max"`
		Region   string  `json:"region" firestore:"region"`
	}

	Settings struct {
		Stores   Store    `json:"stores,omitempty" firestore:"stores,omitempty"`
		Webhooks []string `json:"webhooks,omitempty" firestore:"webhooks,omitempty"`
	}

	Session struct {
		Status    string     `json:"status" firestore:"status"`
		Instances []Instance `json:"instances,omitempty" firestore:"instances,omitempty"`
		Tasks     int        `json:"tasks,omitempty" firestore:"tasks,omitempty"`
	}

	Billing struct {
		CC_Number string `json:"ccnumber" firestore:"cc_number"`
		Month     string `json:"month" firestore:"month"`
		Year      string `json:"year" firestore:"year"`
		CVV       string `json:"cvv" firestore:"cvv"`
	}

	Shipping struct {
		Last     string `json:"last" firestore:"last"`
		Address1 string `json:"address1" firestore:"address1"`
		Address2 string `json:"address2" firestore:"address2"`
		City     string `json:"city" firestore:"city"`
		Email    string `json:"email" firestore:"email"`
		Country  string `json:"country" firestore:"country"`
		State    string `json:"state" firestore:"state"`
		Zip      string `json:"zip"  firestore:"zip"`
	}

	DataResp struct {
		Checkouts []Product  `json:"checkouts,omitempty"`
		Instances []Instance `json:"instances,omitempty"`
		Settings  *Settings  `json:"settings,omitempty"`
		Session   *Session   `json:"session,omitempty"`
		Billing   []Billing  `json:"billing,omitempty"`
		Shipping  []Shipping `json:"shipping,omitempty"`
		User      User       `json:"user"`
		Whitelist []string   `json:"whitelist,omitempty"`
	}

	Notification struct {
		Title     string    `json:"title,omitempty" firestore:"title,omitempty"`
		Text      string    `json:"text,omitempty" firestore:"text,omitempty"`
		CreatedAt time.Time `json:"created_at" firestore:"created_at"`
		Read      *bool     `json:"read,omitempty" firestore:"read,omitempty"`
		Global    bool      `json:"global,omitempty" firestore:"global,omitempty"`
		ID        string    `json:"id" firestore:"-"`
	}
)

func (s *Store) UnmarshalJSON(data []byte) error {
	var settings map[modules.Site]bool

	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	for _, site := range modules.Sites {
		if settings[site] {
			s.Add(site)
		}
	}

	return nil
}

func (s Store) MarshalJSON() ([]byte, error) {
	settings := make(map[modules.Site]bool)
	for _, site := range modules.Sites {
		settings[site] = s.IsEnabled(site)
	}

	return json.Marshal(settings)
}
