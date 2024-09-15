package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DomainRecord struct {
	Domain     string
	TxtRecords []string
}

// Global variables
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

	// Get the collection
	collection = client.Database("recordlookup").Collection("domains")
	if collection == nil {
		log.Fatal("Failed to get collection 'domains'")
	}
}

// Set up the routes for the application
func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/search", searchKeyword)
	http.HandleFunc("/list", getAllDomains)

	// Start the server on port 8080
	http.ListenAndServe(":8080", nil)
}

// upload and process the bulk domains file
func uploadFile(w http.ResponseWriter, r *http.Request) {

	const grMax = 100
	var wg sync.WaitGroup
	ch := make(chan int, grMax)

	//Set a maximum amount of memory to be used when parsing the request body
	r.ParseMultipartForm(10 << 20) // 10 << 20 equivalent to 10mb

	// Get the file from the request
	file, handler, err := r.FormFile("domainsFile")

	if err != nil {
		fmt.Println("Error retriving file from the request")
		fmt.Println(err)
		return
	}

	// Close the file after were done with it
	defer file.Close()

	fmt.Printf("Uploaded file: %v\n", handler.Filename)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Creating a slice to store the domains
	var domains []string
	for scanner.Scan() {
		// Trim the line to remove whitespace
		domain := strings.TrimSpace(scanner.Text())
		domains = append(domains, domain)
	}

	var i int

	// Process each line in the file
	for _, domain := range domains {
		// If the line is not empty then we process the domain
		if domain != "" {
			wg.Add(1)
			ch <- 1
			percent := int(float64(i) / float64(len(domains)) * 100)
			fmt.Printf("Processing domain: %v (%d%%) \n", domain, percent)

			go func() {
				defer func() { wg.Done(); <-ch }()
				processDomain(domain)
			}()
		}

		i++
	}

	wg.Wait()

	// Check if there were any errors while reading the file
	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File processed successfully")
}

// This function takes a domain name, looks up its TXT records, and stores them in the database if they exist.
func processDomain(domain string) {
	// Look up the TXT records for the domain
	// txtRecords, err := lookupTXTWithAPI(domain)
	// txtRecords, err := lookupTXTWithMiekg(domain)
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		fmt.Printf("Error processing TXT records for %v: %v\n", domain, err)
		return
	}

	// Quiting the function if there are no TXT records
	if len(txtRecords) == 0 {
		return
	}

	// Creating a DomainRecord object to store in the database
	record := DomainRecord{Domain: domain, TxtRecords: txtRecords}

	// Inserting the record into the database
	_, err = collection.InsertOne(context.TODO(), record)

	if err != nil {
		fmt.Printf("Error storing record for %s: %v\n", domain, err)
	}
}

// Search the database for the given keyword
func searchKeyword(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	if keyword == "" {
		http.Error(w, "Keyword is required", http.StatusBadRequest)
		return
	}

	// a filter variable to search the database
	filter := bson.M{"txtrecords": bson.M{"$regex": keyword}}

	// Execute the search
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Store the results in a slice of strings
	var domains []string

	// Interate over the results and decoding each record
	for cursor.Next(context.TODO()) {
		var record DomainRecord
		if err = cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}

		domains = append(domains, record.Domain)
	}

	// Writing the results to a file for download
	filePath := fmt.Sprintf("./results_%d.csv", time.Now().Unix())
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, domain := range domains {
		writer.Write([]string{domain})
	}

	fmt.Fprintf(w, "Results written to file: %s", filePath)
}

// Getting all the processed domains from the database
func getAllDomains(w http.ResponseWriter, r *http.Request) {
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var allDomains []DomainRecord
	for cursor.Next(context.TODO()) {
		var record DomainRecord
		if err = cursor.Decode(&record); err != nil {
			http.Error(w, "Error decoding database record", http.StatusInternalServerError)
			return
		}
		allDomains = append(allDomains, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allDomains)
}

// Downloading the processed file
func downloadFile(w http.ResponseWriter, r *http.Request) {
	//reading the query parameter
	filePath := r.URL.Query().Get("file")

	// finding the file
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// returning the file
	http.ServeFile(w, r, filePath)
}

func main() {
	initDB()
	setupRoutes()
}
