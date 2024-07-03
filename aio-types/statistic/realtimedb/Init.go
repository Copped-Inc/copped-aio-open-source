package realtimedb

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"strings"
)

var Database Struct

func Init(server string) error {

	sa := option.WithCredentialsJSON([]byte(Credentials))
	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		DatabaseURL: URL,
	}, sa)
	if err != nil {
		return err
	}

	realtimeclient, err := app.Database(context.Background())
	if err != nil {
		return err
	}

	Database = Struct{
		Server:   strings.ReplaceAll(strings.Split(server, "//")[1], ".", "-"),
		Database: realtimeclient,
		Context:  context.Background(),
	}

	return err

}

type Struct struct {
	Server   string
	Database *db.Client
	Context  context.Context
}
