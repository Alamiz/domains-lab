package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
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

// Global constants
const baseDir = "./results"

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

	// Get the collection
	collection = db.Collection("domains")
	if collection == nil {
		log.Fatal("Failed to get collection 'domains'")
	}

	// Ensure Indexes
	createIndexes(collection)
}

// Create indexes for the 'domains' collection
func createIndexes(collection *mongo.Collection) {

	// Create an index for the 'domain' and 'txtRecords' fields
	indexes := []mongo.IndexModel{
		{Keys: bson.M{"domain": 1}},
		{Keys: bson.M{"txtrecords": "text"}},
	}

	// Ensure both indexes are created
	_, err := collection.Indexes().CreateMany(context.TODO(), indexes)
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

	// Flusher ensures that the response can be written in chunks
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

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

	const grMax = 200           // MAX number of goroutines running simultaneously
	var wg sync.WaitGroup       // WaitGroup to wait for all goroutines to finish
	ch := make(chan int, grMax) // Channel to control goroutine concurrency
	var i int32                 // Counter for the number of processed domains
	var domainsLength int32     // Length of the total domains in the file

	// Create a scanner to read the file line by line
	lengthScanner := bufio.NewScanner(file)
	for lengthScanner.Scan() {
		domainsLength++
	}

	// send the cursor back to the start of the file
	file.Seek(0, 0)

	// Generate a unique file name by appending the current timestamp
	timestamp := time.Now().Unix()
	uniqueFileName := fmt.Sprintf("%s_%d", handler.Filename, timestamp)

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Process each line in the file
	for scanner.Scan() {
		// Trim the line to remove whitespace
		domain := strings.TrimSpace(scanner.Text())

		// If the line is not empty then we process the domain
		if domain == "" {
			return
		}

		// Adding 1 to the wait gorup counter and sending a value to the channel buffer to start the goroutine
		wg.Add(1)
		ch <- 1

		// Launch goroutine to process the domain
		go func() {
			defer func() {
				wg.Done()
				<-ch
			}()
			processDomain(domain, uniqueFileName)
		}()

		percent := int(float64(i) / float64(domainsLength) * 100)

		fmt.Printf("Processing domain: %v (%d%%) \n", domain, percent)
		fmt.Fprintf(w, "%d\n", percent)
		flusher.Flush() // Send the data to the client after the domain is done processing
		i++
	}

	// Check if there were any errors while reading the file
	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	wg.Wait()
	fmt.Println("File processed successfully")
	fmt.Fprintf(w, "File processed successfully")
	fmt.Fprintf(w, "100\n")
	flusher.Flush() // Final flush to ensure the last chunk is sent
}

// This function takes a domain name, looks up its TXT records, and stores them in the database if they exist.
func processDomain(domain string, fileName string) {
	const (
		maxRetries = 3
		retryDelay = 200 * time.Millisecond
	)

	var err error
	var txtRecords []string

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Look up the TXT records for the domain
		txtRecords, err = net.LookupTXT(domain)

		if err != nil {
			if strings.HasSuffix(err.Error(), "no such host") {
				fmt.Printf("No such host: %v\n", domain)
				return
			}

			fmt.Printf("Error processing TXT records for %v (attempt %d/%d): %v\n", domain, attempt+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		break
	}

	if err != nil {
		fmt.Printf("Failed to process TXT records for %v after %d attempts\n", domain, maxRetries)
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
	filter := bson.M{"$text": bson.M{"$search": keyword}}

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
	filePath := fmt.Sprintf("%v/results_%d.csv", baseDir, time.Now().Unix())
	fmt.Println("File path:", filePath)
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
	filePath, err := getFilePath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = serveFile(w, r, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get the file path and check if it exists
func getFilePath(r *http.Request) (string, error) {
	filePath := r.URL.Query().Get("file")
	if !strings.HasPrefix(filePath, baseDir) {
		return "", errors.New("invalid file path")
	}
	return filePath, nil
}

// Send the file through the http response
func serveFile(w http.ResponseWriter, r *http.Request, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	http.ServeFile(w, r, filePath)
	return nil
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
