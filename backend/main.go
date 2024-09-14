package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DomainRecord struct {
	Domain     string
	TxtRecords []string
}

var collection *mongo.Collection

// Initialize MongoDB
func initDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get MongoDB URI from environment
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}

	// Ensure disconnection when done
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v\n", err)
		}
	}()

	// Get the collection
	collection := client.Database("recordlookup").Collection("domains")
	if collection == nil {
		log.Fatal("Failed to get collection 'domains'")
	}
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/search", searchKeyword)
	http.ListenAndServe(":8080", nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File uploaded")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Downloading file")
}

func searchKeyword(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Searching for the keyword ...")
}

func main() {
	initDB()
	setupRoutes()
}
