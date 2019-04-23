//Package fdb handles firebase DB transactions
package fdb

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	ctx context.Context
)

const (
	playersCollection = "players"
	winnersCollection = "winners"
	sessionCollection = "session"
)

// Player represents the struct of player in players db
type Player struct {
	Email   string
	Key     string
	Session int
	Amount  int
}

// Fdb is the struct to communicate with the Lottery App Firebase DB
type Fdb struct {
	client  *firestore.Client
	Players []Player
}

//NewFdb start a new fdb connection
func NewFdb(key, project string) (*Fdb, error) {
	ctx = context.Background()

	opt := option.WithCredentialsFile(key)
	fdb := new(Fdb)

	client, err := firestore.NewClient(ctx, project, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client: %v", err)
	}
	fdb.client = client
	fdb.Players = make([]Player, 0)

	return fdb, nil
}

// CloseFdb will close the Firebase DB connection
func (fdb *Fdb) CloseFdb() {
	fdb.client.Close()
}

// Find more examples here: https://cloud.google.com/firestore/docs/quickstart-servers

// AddData --
func addData(client *firestore.Client) {
	// _, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
	// 	"first": "Ada",
	// 	"last":  "Lovelace",
	// 	"born":  1815,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed adding alovelace: %v", err)
	// }
}

// readData read data from firebase
func (fdb *Fdb) readData(collection string) []interface{} {
	iter := fdb.client.Collection(collection).Documents(ctx)
	data := make([]interface{}, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Failed to iterate: %v", err)
			continue
		}
		data = append(data, doc.Data())
	}
	return data
}

//GetSession return the currrent session number
func (fdb *Fdb) GetSession() int {
	return 0
}

//GetPlayers returns a list of players in the current session
//all: true, return all players
func (fdb *Fdb) GetPlayers(all bool) int {
	if all {
		players := fdb.readData(playersCollection)
		for i, p := range players {
			fmt.Printf("%v => %v\n", i, p)
		}
		return len(players)
	}
	return 0
}
