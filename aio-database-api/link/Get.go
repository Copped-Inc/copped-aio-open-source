package link

import "database-api/database"

func Get(id string) (*Link, error) {
	found, err := database.GetDatabase().Collection("link").Doc(id).Get(database.GetContext())
	if err != nil {
		return &Link{}, err
	}

	var l *Link = &Link{ID: id}
	return l, found.DataTo(l)
}
