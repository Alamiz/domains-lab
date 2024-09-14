package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("domainsFile")

	if err != nil {
		fmt.Println("Error retriving file from the request")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Printf("Uploaded file: %v\n", handler.Filename)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			go processDomain(domain)
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "File processed successfully")
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Downloading file")
}

func searchKeyword(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	if keyword == "" {
		http.Error(w, "Keyword is required", http.StatusBadRequest)
		return
	}

	filter := bson.D{{"txtrecords", bson.M{"$regex": keyword}}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var domains []string
	for cursor.Next(context.TODO()) {
		var record DomainRecord
		if err = cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}

		domains = append(domains, record.Domain)
	}
}

func processDomain(domain string) {
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		fmt.Printf("Error processing TXT records for %v: %v\n", domain, err)
		return
	}

	if len(txtRecords) > 0 {
		record := DomainRecord{Domain: domain, TxtRecords: txtRecords}
		_, err := collection.InsertOne(context.TODO(), record)
		if err != nil {
			fmt.Printf("Error storing record for %s: %v\n", domain, err)
		}
	}
}

func main() {
	initDB()
	setupRoutes()
}
