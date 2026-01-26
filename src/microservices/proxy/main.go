package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const envMoviesMigrationPercent = "MOVIES_MIGRATION_PERCENT"
const envMonolithUrl = "MONOLITH_URL"
const envMoviesServiceUrl = "MOVIES_SERVICE_URL"
const envEventsServiceUrl = "EVENTS_SERVICE_URL"
const envGradualMigration = "GRADUAL_MIGRATION"

var moviesMigrationPercent int
var monolithUrl string
var moviesServiceUrl string
var eventsServiceUrl string
var gradualMigration bool

var counter *BinaryCounter

func main() {
	getEvns()

	counter = createBinaryCounter(moviesMigrationPercent, 100-moviesMigrationPercent)

	// Set up HTTP routes
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/api/movies", handleMovies)
	http.HandleFunc("/api/events", handleEvents)
	http.HandleFunc("/api/", handleOther)

	// Start server
	port, _ := lookupEnv("PORT")
	if port == "" {
		port = "8081" // Note: Using a different port than the monolith
	}

	log.Printf("Starting movies microservice on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Handlers ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	if counter.next() == 1 {
		forwardRequest(w, r, moviesServiceUrl)
	} else {
		forwardRequest(w, r, monolithUrl)
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, eventsServiceUrl)
}

func handleOther(w http.ResponseWriter, r *http.Request) {
	forwardRequest(w, r, monolithUrl)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// BinaryCounter ///////////////////////////////////////////////////////////////////////////////////////////////////////

/*
BinaryCounter holds state of two counters.
*/
type BinaryCounter struct {
	firstMax      int
	secondMax     int
	firstCounter  int
	secondCounter int
}

func createBinaryCounter(firstMax int, secondMax int) *BinaryCounter {
	return &BinaryCounter{
		firstMax:      firstMax,
		secondMax:     secondMax,
		firstCounter:  0,
		secondCounter: 0}
}

func (bc *BinaryCounter) next() int {
	if bc.firstCounter == 0 && bc.secondCounter == 0 {
		bc.firstCounter = bc.firstMax
		bc.secondCounter = bc.secondMax
	}

	if bc.firstCounter != 0 && bc.secondCounter != 0 {
		var next = rand.Intn(2) + 1
		switch next {
		case 1:
			bc.firstCounter--
			break
		case 2:
			bc.secondCounter--
			break
		}

		return next
	} else if bc.firstCounter != 0 {
		bc.firstCounter--
		return 1
	} else { // bc.secondCounter != 0
		bc.secondCounter--
		return 2
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Helpers ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func getEvns() {
	val, exists := lookupEnv(envMonolithUrl)
	if exists {
		monolithUrl = val
	} else {
		panic("Mandatory environment variable not found: " + envMonolithUrl)
	}

	env, exists := lookupEnv(envGradualMigration)
	if exists {
		gradualMigration = parseBool(env)
	} else {
		gradualMigration = false
	}

	if gradualMigration {
		moviesServiceUrl, exists = lookupEnv(envMoviesServiceUrl)
		if !exists {
			panic("Mandatory environment variable not found: " + envMoviesServiceUrl)
		}

		eventsServiceUrl, exists = lookupEnv(envEventsServiceUrl)
		if !exists {
			// todo:
		}

		val, exists = os.LookupEnv(envMoviesMigrationPercent)
		if exists {
			moviesMigrationPercent = parseInt(val)
		} else {
			panic("Mandatory environment variable not found: " + envMoviesMigrationPercent)
		}
	}
}

func lookupEnv(envName string) (string, bool) {
	val, exists := os.LookupEnv(envName)

	if exists {
		log.Printf("Environment variable (%s): %s", envName, val)
	} else {
		log.Printf("Environment variable (%s): not found", envName)
	}

	return val, exists
}

func parseBool(env string) bool {
	val, err := strconv.ParseBool(env)

	if err != nil {
		panic(err)
	}

	return val
}

func parseInt(env string) int {
	val, err := strconv.ParseInt(env, 10, 32)

	if err != nil {
		panic(err)
	}

	return int(val)
}

func forwardRequest(w http.ResponseWriter, r *http.Request, forwardedHost string) {
	// Log the incoming request for debugging
	log.Printf("Received request: %s %s from %s", r.Method, r.URL, r.RemoteAddr)

	// Construct the target URL - in a forward proxy, the client specifies the full URL
	var targetURL string
	if forwardedHost == "" {
		if r.URL.Scheme == "" {
			// If no scheme provided, default to HTTP
			targetURL = "http://" + r.Host + r.URL.Path
		} else {
			targetURL = r.URL.String()
		}
	} else {
		targetURL = forwardedHost + r.URL.Path
	}

	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	// Create the outbound proxy request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers from the original request, excluding hop-by-hop headers
	for name, values := range r.Header {
		// Skip headers that shouldn't be forwarded
		if name == "Proxy-Connection" || name == "Connection" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// Add X-Forwarded-For header to track the original client
	if clientIP := r.RemoteAddr; clientIP != "" {
		proxyReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// Execute the request using our custom transport
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		log.Printf("Error forwarding request: %v", err)
		http.Error(w, "Failed to reach destination server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers back to the client
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// Forward the status code
	w.WriteHeader(resp.StatusCode)

	// Stream the response body back to the client
	written, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
		return
	}

	log.Printf("Proxied %s %s - Status: %d, Bytes: %d", r.Method, targetURL, resp.StatusCode, written)
}
