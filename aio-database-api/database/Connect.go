package database

import (
	"context"
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var Database Struct

func Connect() error {

	opt := option.WithCredentialsJSON([]byte(realtimedb.Credentials))

	client, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	database, err := client.Firestore(context.Background())
	if err != nil {
		return err
	}

	Database = Struct{
		Database: database,
		Context:  context.Background(),
	}

	return err

}

type Struct struct {
	Database *firestore.Client
	Context  context.Context
}
