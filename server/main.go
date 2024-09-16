package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DomainRecord struct {
	Domain     string
	TxtRecords []string
	FileName   string
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

	db := client.Database("recordlookup")
	// Ensure Indexes
	createIndexes(db)

	// Get the collection
	collection = client.Database("recordlookup").Collection("domains")
	if collection == nil {
		log.Fatal("Failed to get collection 'domains'")
	}
}

// Create indexes for the 'domains' collection
func createIndexes(db *mongo.Database) {
	collection := db.Collection("domains")

	// Create an index for the 'domain' field
	domainIndex := mongo.IndexModel{
		Keys: bson.M{"domain": 1}, // Indexing the 'domain' field in ascending order
	}

	// Create an index for the 'txtRecords' field
	txtRecordsIndex := mongo.IndexModel{
		Keys: bson.M{"txtRecords": "text"}, // Full-text search index on 'txtRecords'
	}

	// Ensure both indexes are created
	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{domainIndex, txtRecordsIndex})
	if err != nil {
		log.Fatalf("Error creating indexes: %v\n", err)
	}

	fmt.Println("Indexes created successfully")
}

// CORS Middleware function
func enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow requests from the client origin

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

		// Allow specified HTTP methods

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specified headers

		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		// Handle OPTIONS requests (preflight request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue with the next handler

		next.ServeHTTP(w, r)

	})
}

// Set up the routes for the application
func setupRoutes() {
	router := mux.NewRouter()

	// Enables CORS
	router.Use(enableCORS)

	// Routes
	router.HandleFunc("/upload", uploadFile).Methods("POST")
	router.HandleFunc("/download", downloadFile).Methods("GET")
	router.HandleFunc("/search", searchKeyword).Methods("GET")
	router.HandleFunc("/list", getAllDomains).Methods("GET")

	// Start the server on port 8080
	http.ListenAndServe(":8080", router)
}

// upload and process the bulk domains file
func uploadFile(w http.ResponseWriter, r *http.Request) {

	// Setting headers for streaming
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "keep-alive")

	//Set a maximum amount of memory to be used when parsing the request body
	r.ParseMultipartForm(10 << 20) // 10 << 20 equivalent to 10mb

	// Get the file from the request
	file, handler, err := r.FormFile("domainsFile")

	if err != nil {
		fmt.Println("Error retriving file from the request")
		fmt.Println(err)
		return
	}

	if !validateFileType(handler) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
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

	// Flusher ensures that the response can be written in chunks
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	const grMax = 100           // MAX number of goroutines running simultaneously
	var wg sync.WaitGroup       // WaitGroup to wait for all goroutines to finish
	ch := make(chan int, grMax) // Channel to control goroutine concurrency
	var i int32                 // Use atomic counter for thread-safe increments

	// Generate a unique file name by appending the current timestamp
	timestamp := time.Now().Unix()
	uniqueFileName := fmt.Sprintf("%s_%d", handler.Filename, timestamp)

	// Process each line in the file
	for _, domain := range domains {
		// If the line is not empty then we process the domain
		if domain != "" {
			wg.Add(1)
			ch <- 1
			percent := int(float64(i) / float64(len(domains)) * 100)
			fmt.Printf("Processing domain: %v (%d%%) \n", domain, percent)

			fmt.Fprintf(w, "%d\n", percent)
			flusher.Flush() // Send the data to the client immediately

			// Launch goroutine to process the domain
			go func() {
				defer func() { wg.Done(); <-ch }()
				processDomain(domain, uniqueFileName)
			}()
		}

		i++
	}

	wg.Wait()
	fmt.Println("File processed successfully")
	fmt.Fprintf(w, "100\n")
	flusher.Flush() // Final flush to ensure the last chunk is sent

	// Check if there were any errors while reading the file
	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File processed successfully")
}

// This function takes a domain name, looks up its TXT records, and stores them in the database if they exist.
func processDomain(domain string, fileName string) {
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
	record := DomainRecord{Domain: domain, TxtRecords: txtRecords, FileName: fileName}

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
	var domains []DomainRecord

	// Interate over the results and decoding each record
	for cursor.Next(context.TODO()) {
		var record DomainRecord
		if err = cursor.Decode(&record); err != nil {
			log.Fatal(err)
		}

		domains = append(domains, record)
	}

	// If there are no results, return an error
	if len(domains) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
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
		writer.Write([]string{domain.Domain, domain.FileName})
	}

	// Send the filepath as a response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"filepath\":\"%s\"}", filePath)
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

// Validate file type
func validateFileType(header *multipart.FileHeader) bool {
	fileName := header.Filename
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".csv", ".txt":
		return true
	default:
		return false
	}
}

func main() {
	initDB()
	setupRoutes()
}
