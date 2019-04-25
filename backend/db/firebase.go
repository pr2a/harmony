//Package fdb handles firebase DB transactions
package fdb

import (
	"context"
	"fmt"
	"log"
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

func (p *Player) String() string {
	return fmt.Sprintf("player:%s/%s (%v/%v)", p.Address, p.Email, p.Amount, p.Balance)
}

// Winner of the lottery
type Winner struct {
	Session int64
	Amount  int64
	Address string
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
		onePlayer := new(Player)
		onePlayer.Address = p
		onePlayer.Balance = new(big.Int)
		onePlayer.Balance.SetString(players.Balances[i], 10)
		dbPlayers = append(dbPlayers, onePlayer)
	}
	return dbPlayers
}

// Find more examples here: https://cloud.google.com/firestore/docs/quickstart-servers

// UpdateSession will set is_current to 'false'
func (fdb *Fdb) UpdateSession() {
	q := fdb.client.Collection(sessionCollection).Select().Where("is_current", "==", true)
	iter := q.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	_, err = doc.Ref.Set(ctx, map[string]interface{}{
		"is_current": false,
	}, firestore.MergeAll)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
}

// AddSession add new session id to session collection
func (fdb *Fdb) AddSession(id int64) {
	_, _, err := fdb.client.Collection(sessionCollection).Add(ctx, map[string]interface{}{
		"deadline":   time.Now(),
		"is_current": true,
		"session_id": id,
	})
	if err != nil {
		log.Fatalf("Failed adding a new session: %v", err)
	}
}

// AddWinner add a new document in winner collection
func (fdb *Fdb) AddWinner(win *Winner) {
	_, _, err := fdb.client.Collection(winnersCollection).Add(ctx, map[string]interface{}{
		"amount":     win.Amount,
		"session_id": win.Session,
		"address":    win.Address,
	})
	if err != nil {
		log.Fatalf("Failed adding a new session: %v", err)
	}
}

// readData read data from firebase
func (fdb *Fdb) readData(collection string) []interface{} {
	iter := fdb.client.Collection(collection).Documents(ctx)
	defer iter.Stop()

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
		iter = fdb.client.Collection(sessionCollection).Where("is_current", "==", true).Documents(ctx)
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
		one.Current, ok = data["is_current"].(bool)
		if !ok {
			fmt.Printf("Failed to convert \"is_current\"\n")
			continue
		}
		one.ID, ok = data["session_id"].(int64)
		if !ok {
			fmt.Printf("Failed to convert \"session_id\"\n")
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
