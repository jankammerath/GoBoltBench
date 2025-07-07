package main

import (
	"log"
	"math/rand"
	"os"
	"strings"

	"go.etcd.io/bbolt"
)

type Message struct {
	UserName    string `json:"userName"`
	MessageText string `json:"messageText"`
}

func getUsers() []string {
	file, err := os.ReadFile("users.txt")
	if err != nil {
		log.Fatalf("Failed to read users file: %v", err)
	}
	return strings.Split(string(file), "\n")
}

func getLines() []string {
	file, err := os.ReadFile("testfile.txt")
	if err != nil {
		log.Fatalf("Failed to read test file: %v", err)
	}
	return strings.Split(string(file), "\n")
}

func getRandomUser(users []string) string {
	if len(users) == 0 {
		log.Fatal("No users available")
	}

	// return a random user from the list
	index := rand.Intn(len(users))
	return users[index]
}

func main() {
	const workerCount = 16
	const dbFile = "output.db"
	const bucketName = "messages"

	// Remove old database file
	if _, err := os.Stat(dbFile); err == nil {
		if err := os.Remove(dbFile); err != nil {
			log.Fatalf("Failed to remove old database file: %v", err)
		}
	}

	// Open the bbolt database.
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create a bucket.
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	users := getUsers()
	lines := getLines()
	if len(users) == 0 || len(lines) == 0 {
		log.Fatal("Users or lines are empty, please check the input files.")
	}

}
