//Package fdb handles firebase DB transactions
package fdb

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/harmony-one/demo-apps/backend/client"
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
	Email    string
	RegEmail bool
	PrivKey  string
	Address  string
	Session  int64
	Notified bool
	Amount   int
	Balance  *big.Int
}

// Session represents the struct of the session in session collection
type Session struct {
	Deadline time.Time
	Current  bool
	ID       int64
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

// NewPlayer convert db struct from rest API to DB player
func NewPlayer(players *restclient.Player) []*Player {
	dbPlayers := make([]*Player, 0)
	for i, p := range players.Players {
		n := new(big.Int)
		n, ok := n.SetString(players.Balances[i], 10)
		if !ok {
			fmt.Printf("SetString Error")
			continue
		}
		onePlayer := Player{
			Address: p,
			Balance: n,
		}
		dbPlayers = append(dbPlayers, &onePlayer)
	}
	return dbPlayers
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
func (fdb *Fdb) GetSession(all bool) []*Session {
	var iter *firestore.DocumentIterator
	var ok bool

	if all {
		iter = fdb.client.Collection(sessionCollection).Documents(ctx)
	} else {
		iter = fdb.client.Collection(sessionCollection).Where("current", "==", true).Documents(ctx)
	}
	sessions := make([]*Session, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Failed to iterate: %v", err)
			continue
		}
		data := doc.Data()
		one := new(Session)

		one.Deadline, ok = data["deadline"].(time.Time)
		if !ok {
			fmt.Printf("Failed to convert \"deadline\": %v\n", data["deadline"])
			continue
		}
		one.Current, ok = data["current"].(bool)
		if !ok {
			fmt.Printf("Failed to convert \"current\"\n")
			continue
		}
		one.ID, ok = data["id"].(int64)
		if !ok {
			fmt.Printf("Failed to convert \"id\"\n")
			continue
		}
		sessions = append(sessions, one)
	}
	return sessions
}

//GetPlayers returns a list of players in the current session
//all: true, return all players
func (fdb *Fdb) GetPlayers(key, op string, value interface{}) []*Player {
	var iter *firestore.DocumentIterator
	var ok bool

	if key == "" {
		iter = fdb.client.Collection(playersCollection).Documents(ctx)
	} else {
		iter = fdb.client.Collection(playersCollection).Where(key, op, value).Documents(ctx)
	}

	players := make([]*Player, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("Failed to iterate: %v", err)
			continue
		}
		data := doc.Data()
		one := new(Player)

		one.Email, ok = data["email"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"email\": %v\n")
			continue
		}
		one.RegEmail, ok = data["keys_notified"].(bool)
		if !ok {
			fmt.Printf("Failed to convert \"keys_notified\"\n")
			continue
		}
		one.PrivKey, ok = data["private_key"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"private_key\"\n")
			continue
		}
		one.Address, ok = data["address"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"address\"\n")
			continue
		}
		one.Session, ok = data["session_id"].(int64)
		if !ok {
			fmt.Printf("Failed to convert \"session_id\"\n")
			continue
		}
		one.Notified, ok = data["result_notified"].(bool)
		if !ok {
			fmt.Printf("Failed to convert \"result_notified\"\n")
			continue
		}

		players = append(players, one)
	}
	return players
}
