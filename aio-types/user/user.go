package user

import (
	"time"

	"github.com/Copped-Inc/aio-types/subscriptions"
)

type User struct {
	Name          string       `json:"name" firestore:"name"`
	Email         string       `json:"email" firestore:"email"`
	ID            string       `json:"id" firestore:"id"`
	Picture       string       `json:"picture" firestore:"picture"`
	InstanceLimit int          `json:"instance_limit" firestore:"instance_limit"`
	Code          string       `json:"code,omitempty" firestore:"code,omitempty"`
	CodeExpire    time.Time    `json:"code_expire,omitempty" firestore:"code_expire,omitempty"`
	Subscription  Subscription `json:"subscription" firestore:"subscription"`
}

type Subscription struct {
	Plan            subscriptions.Plan `json:"plan" firestore:"plan"`
	Customer_ID     string             `json:"customer_id,omitempty" firestore:"customer_id,omitempty"`
	Subscription_ID string             `json:"subscription_id,omitempty" firestore:"subscription_id,omitempty"`
	State           State              `json:"state" firestore:"state"`
}

type State int

const (
	Active State = iota + 1
	Pending
	Disabled
	Expired
)
