//go:build !mem
// +build !mem

package data

import (
	"fmt"

	"ytsruh.com/saas/data/model"
	"ytsruh.com/saas/data/mongo"
)

func (db *DB) Open(driverName, dataSourceName string) error {
	fmt.Println("inside mongo")
	conn, err := model.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	//  for mongo, we need to copy the connection session at each requests
	// this is done in our api's ServeHTTP
	db.Users = &mongo.Users{}
	db.Webhooks = &mongo.Webhooks{}

	db.Connection = conn

	db.DatabaseName = "gosaas"
	db.CopySession = true
	return nil
}
