package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

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

// worker processes a slice of lines and stores them as messages in the database
func worker(db *bbolt.DB, bucketName string, lines []string, users []string, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for i, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Create a message with random user and the line text
		message := Message{
			UserName:    getRandomUser(users),
			MessageText: line,
		}

		// Marshal message to JSON
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("Worker %d: Failed to marshal message: %v", workerID, err)
			continue
		}

		// Store in database
		err = db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				return fmt.Errorf("bucket %s not found", bucketName)
			}

			// Generate a unique key for each message
			key := fmt.Sprintf("worker_%d_msg_%d", workerID, i)
			return bucket.Put([]byte(key), messageJSON)
		})

		if err != nil {
			log.Printf("Worker %d: Failed to store message: %v", workerID, err)
		}
	}

	log.Printf("Worker %d completed processing %d lines", workerID, len(lines))
}

func main() {
	startTime := time.Now()

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

	// Filter out empty lines
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	lines = filteredLines

	log.Printf("Processing %d lines with %d workers", len(lines), workerCount)

	// Calculate lines per worker
	linesPerWorker := len(lines) / workerCount
	remainder := len(lines) % workerCount

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		// Calculate start and end indices for this worker
		start := i * linesPerWorker
		end := start + linesPerWorker

		// Distribute remainder lines to the first few workers
		if i < remainder {
			end++
		}

		// Adjust start for workers that get extra lines
		if i > 0 && i <= remainder {
			start += i
		} else if i > remainder {
			start += remainder
		}

		// Ensure we don't go out of bounds
		if end > len(lines) {
			end = len(lines)
		}

		// Skip if this worker has no lines to process
		if start >= len(lines) {
			wg.Done()
			continue
		}

		workerLines := lines[start:end]
		go worker(db, bucketName, workerLines, users, &wg, i)
	}

	// Wait for all workers to complete
	wg.Wait()
	log.Println("All workers completed successfully")

	totalTime := time.Since(startTime)
	log.Printf("Total processing time: %v", totalTime)
}
