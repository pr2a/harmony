package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	TestPlayers = "test_players"
)

// Find more examples here: https://cloud.google.com/firestore/docs/quickstart-servers

// AddData --
func AddData(client *firestore.Client) {
	// _, _, err = client.Collection("users").Add(ctx, map[string]interface{}{
	// 	"first": "Ada",
	// 	"last":  "Lovelace",
	// 	"born":  1815,
	// })
	// if err != nil {
	// 	log.Fatalf("Failed adding alovelace: %v", err)
	// }
}

// ReadData --
func ReadData(client *firestore.Client, ctx context.Context) {
	iter := client.Collection(TestPlayers).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}

}

func main() {
	projectID := "benchmark-209420"
	// Get a Firestore client.
	ctx := context.Background()
	opt := option.WithCredentialsFile("./keys/benchmark_account_key.json")
	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Close client when done.
	defer client.Close()

	ReadData(client, ctx)
}
