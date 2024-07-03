package link

import "github.com/Copped-Inc/aio-types/subscriptions"

type Link struct {
	ID            string             `json:"id" firestore:"-"`
	Plan          subscriptions.Plan `json:"plan" firestore:"plan"`
	Stock         int                `json:"stock" firestore:"stock"`
	InstanceLimit int                `json:"instance_limit" firestore:"instance_limit"`
}
