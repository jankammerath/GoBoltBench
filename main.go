package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/mem"
	"go.etcd.io/bbolt"
)

type Message struct {
	UserName    string `json:"userName"`
	MessageText string `json:"messageText"`
	Hash        string `json:"hash"`
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

	// Use strings.Fields instead of Split to avoid empty strings
	lines := strings.Fields(string(file))
	return lines
}

func sha256Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func worker(db *bbolt.DB, bucketName string, lines []string, users []string, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	// Pre-generate random users for this worker's messages
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
	randomUsers := make([]string, len(lines))
	for i := 0; i < len(lines); i++ {
		randomUsers[i] = users[rng.Intn(len(users))]
	}

	const batchSize = 100
	var batch []Message
	var batchKeys []string

	flushBatch := func() {
		if len(batch) == 0 {
			return
		}

		err := db.Update(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				return fmt.Errorf("bucket %s not found", bucketName)
			}

			for i, message := range batch {
				messageJSON, err := json.Marshal(message)
				if err != nil {
					log.Printf("Worker %d: Failed to marshal message: %v", workerID, err)
					continue
				}

				if err := bucket.Put([]byte(batchKeys[i]), messageJSON); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			log.Printf("Worker %d: Failed to store batch: %v", workerID, err)
		}

		batch = batch[:0]
		batchKeys = batchKeys[:0]
	}

	messageIndex := 0
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		message := Message{
			UserName:    randomUsers[i],
			MessageText: line,
			// calculate the sha256 hash of the message text
			Hash: sha256Hash(line),
		}

		batch = append(batch, message)
		batchKeys = append(batchKeys, fmt.Sprintf("worker_%d_msg_%d", workerID, messageIndex))
		messageIndex++

		if len(batch) >= batchSize {
			flushBatch()
		}
	}

	// Flush remaining messages
	flushBatch()
}

func getSysMemory() string {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Error getting memory info: %v", err)
	}

	// return in GB
	if v.Total < 1e9 {
		return fmt.Sprintf("%.2f MB", float64(v.Total)/1e6)
	}
	if v.Total < 1e12 {
		return fmt.Sprintf("%.2f GB", float64(v.Total)/1e9)
	}
	if v.Total < 1e15 {
		return fmt.Sprintf("%.2f TB", float64(v.Total)/1e12)
	}
	return fmt.Sprintf("%d", v.Total)
}

func getCpuName() string {
	result := cpuid.CPU.BrandName
	if result == "" {
		// get from /cat/proc/cpuinfo if BrandName is empty
		data, err := os.ReadFile("/proc/cpuinfo")
		if err != nil {
			log.Fatalf("Failed to read /proc/cpuinfo: %v", err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "vendor_id") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) > 1 {
					result = strings.TrimSpace(parts[1])
					break
				}
			}
		}
		if result == "" {
			log.Println("Could not determine CPU name from /proc/cpuinfo, using default")
			result = "Unknown CPU"
		}
	}

	return result
}

func main() {
	bannerText := fmt.Sprintf("GoBoltBench — %s (%s)", getCpuName(), getSysMemory())
	fmt.Println(bannerText)

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

	// Open the bbolt database with optimized options
	dbOptions := &bbolt.Options{
		Timeout:        1 * time.Second,
		NoGrowSync:     true, // Don't sync after growing the database
		NoFreelistSync: true, // Don't sync freelist to disk
		FreelistType:   bbolt.FreelistMapType,
	}

	// Open the bbolt database.
	db, err := bbolt.Open(dbFile, 0600, dbOptions)
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

	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		start := i * len(lines) / workerCount
		end := (i + 1) * len(lines) / workerCount
		if i == workerCount-1 {
			end = len(lines) // Last worker gets remaining lines
		}

		go worker(db, bucketName, lines[start:end], users, &wg, i)
	}
	// Wait for all workers to complete
	wg.Wait()
	log.Println("All workers completed successfully")

	totalTime := time.Since(startTime)
	log.Printf("Total processing time: %v", totalTime)
}
