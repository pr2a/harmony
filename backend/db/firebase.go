//Package fdb handles firebase DB transactions
package fdb

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	restclient "github.com/harmony-one/demo-apps/backend/client"
)

var (
	ctx context.Context
)

const (
	playersCollection = "players"
	winnersCollection = "winners"
	sessionCollection = "session"

	pzPlayersCollection = "pzPlayers"
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
	return fmt.Sprintf("player:%s/%s, session:%v", p.Address, p.Email, p.Session)
}

// PzPlayer represents the struct of puzzle player in the players db
type PzPlayer struct {
	Email   string `firestore:"email"`
	CosID   string `firestore:"cosid"`
	PrivKey string `firestore:"privkey"`
	Address string `firestore:"address"`
	Highest int64  `firestore:"highest"`
	Rewards int64  `firestore:"rewards"`
	Leader  string `firestore:"leader"`
	Port    string `firestore:"port"`
}

func (p *PzPlayer) String() string {
	return fmt.Sprintf("pzPlayer: %s/%s, key: %s/%s, leader: %v", p.Email, p.CosID, p.Address, p.PrivKey, p.Leader)
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

func validateDocsType(docs interface{},
) (slice reflect.Value, elemType reflect.Type, ok bool) {
	var v reflect.Value
	if v = reflect.ValueOf(docs); v.Kind() != reflect.Ptr {
		return
	}
	if v = v.Elem(); v.Kind() != reflect.Slice {
		return
	}
	slice = v
	t := v.Type()
	if t = t.Elem(); t.Kind() != reflect.Ptr {
		return
	}
	if t = t.Elem(); t.Kind() != reflect.Struct {
		return
	}
	elemType = t
	ok = true
	return
}

// QueryFunc is a function that extends the given query.
// It is used as a query filter in Find* methods.
type QueryFunc func(q firestore.Query) firestore.Query

// Find queries the given collection with the given query,
// and appends docs with the documents found.
// docs should be a *[]*DocType, where DocType is a struct.
func (fdb *Fdb) Find(
	ctx context.Context, collection string, queryFunc QueryFunc,
	docs interface{},
) error {
	return fdb.Update(ctx, collection, queryFunc, nil, docs)
}

// Update updates all documents in the given collection that match the given
// query, and appends docs with the documents found.
// The returned documents are at their old values, before updates occurred.
// docs should be a *[]*DocType, where DocType is a struct.
func (fdb *Fdb) Update(
	ctx context.Context, collection string, queryFunc QueryFunc,
	updates []firestore.Update, docs interface{},
) error {
	slice, elemType, ok := validateDocsType(docs)
	if !ok {
		return fmt.Errorf("docs is not a *[]*StructType but a %T", docs)
	}
	query := queryFunc(fdb.client.Collection(collection).Query)
	iter := query.Documents(ctx)
	for fdbDoc, err := iter.Next(); err != iterator.Done; fdbDoc, err = iter.Next() {
		if err != nil {
			return err
		}
		docPtr := reflect.New(elemType)
		if err := fdbDoc.DataTo(docPtr.Interface()); err != nil {
			return err
		}
		if len(updates) > 0 {
			// Prevent update races by requiring that the update time is intact.
			_, err := fdbDoc.Ref.Update(ctx, updates,
				firestore.LastUpdateTime(fdbDoc.UpdateTime))
			if err != nil {
				return err
			}
		}
		slice.Set(reflect.Append(slice, docPtr))
	}
	return nil
}

// FindPzPlayers queries the puzzle player collection with the given query,
// and returns the documents found.
func (fdb *Fdb) FindPzPlayers(
	ctx context.Context, queryFunc QueryFunc,
) (pzPlayers []*PzPlayer, err error) {
	return fdb.UpdatePzPlayers(ctx, queryFunc, nil)
}

// UpdatePzPlayers queries the puzzle player collection with the given query,
// and returns the documents found, and updates them too.
func (fdb *Fdb) UpdatePzPlayers(
	ctx context.Context, queryFunc QueryFunc, updates []firestore.Update,
) (pzPlayers []*PzPlayer, err error) {
	err = fdb.Update(ctx, pzPlayersCollection, queryFunc, updates, &pzPlayers)
	fmt.Printf("(*Fdb).UpdatePzPlayers: pzPlayers=%#v err=%#v\n",
		pzPlayers, err)
	return
}

// FindAccount find the account and leader info from the db
// Assuming the search condition is "key == value"
// Only return the first one account found
func (fdb *Fdb) FindAccount(key string, value string) []*PzPlayer {
	var iter *firestore.DocumentIterator
	var ok bool

	iter = fdb.client.Collection(pzPlayersCollection).Where(key, "==", value).Limit(1).Documents(ctx)

	// We should have only one player returned
	players := make([]*PzPlayer, 0)

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
		one := new(PzPlayer)

		one.Email, ok = data["email"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"email\": %v\n")
			continue
		}
		one.CosID, ok = data["cosid"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"cosid\": %v\n")
			continue
		}
		one.PrivKey, ok = data["privkey"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"privkey\": %v\n")
			continue
		}
		one.Address, ok = data["address"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"address\": %v\n")
			continue
		}
		one.Highest, ok = data["highest"].(int64)
		if !ok {
			fmt.Printf("Failed to convert \"highest\": %v\n")
			continue
		}
		one.Rewards, ok = data["rewards"].(int64)
		if !ok {
			fmt.Printf("Failed to convert \"rewards\": %v\n")
			continue
		}
		one.Leader, ok = data["leader"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"leader\": %v\n")
			continue
		}
		one.Port, ok = data["port"].(string)
		if !ok {
			fmt.Printf("Failed to convert \"port\": %v\n")
			continue
		}

		players = append(players, one)
	}

	return players
}

// RegisterAccount register user account into db
func (fdb *Fdb) RegisterAccount(p *PzPlayer) error {
	_, _, err := fdb.client.Collection(pzPlayersCollection).Add(ctx, map[string]interface{}{
		"email":   p.Email,
		"cosid":   p.CosID,
		"privkey": p.PrivKey,
		"address": p.Address,
		"highest": 0,
		"rewards": 0,
		"leader":  p.Leader,
		"port":    p.Port,
	})

	if err != nil {
		log.Fatalf("Failed adding new account")
	}
	return err
}
