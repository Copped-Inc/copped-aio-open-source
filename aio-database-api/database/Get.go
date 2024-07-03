package database

import (
	"cloud.google.com/go/firestore"
	"context"
)

func GetDatabase() *firestore.Client {

	return Database.Database

}

func GetContext() context.Context {

	return Database.Context

}
