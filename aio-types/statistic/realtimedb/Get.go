package realtimedb

import (
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
)

func GetServer() string {

	return Database.Server

}

func GetDatabase() *db.Client {

	return Database.Database

}

func GetContext() context.Context {

	return Database.Context

}
