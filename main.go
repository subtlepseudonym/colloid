/* Colloid Server
 * subtlepseudonym (subtlepseudonym@gmail.com)
 */

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

type LogEntry struct {
	Id        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"`
	Entry     string    `json:"entry"`
}

const (
	servPort string = ":16874"
	// FIXME: properly set up static ips
	dbAddress string = "postgres://postgres:postgres@192.168.1.3/colloiddb?sslmode=require"
)

var (
	db *sql.DB
)

func simpleJsonResponse(statusStr, msgStr string) []byte {
	return []byte(fmt.Sprintf(`{"status":"%s","msg":"%s"}`, statusStr, msgStr))
}

// TODO: checkout golang's options for serving files (and whether that's easy / a good or effective idea)
// How does golang do with dynamically adding data to the views we're rendering?
// Can always reformat handlers into a REST api and field request from a decoupled front-end
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: serve bootstrap locally (rather than cdn in index.html)
	http.ServeFile(w, r, "assets/index.html")
}

func addLogHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: insert log into db
	// pretty straight forward
	// probably only accepting POST
}

// TODO: add an editLogHandler() ?? maybe updateLogHandler()
// would need to be accessible through the ui for getting logs

func getLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// retrieve all logs
		logs, err := logQuery("SELECT * FROM logs ORDER BY timestamp DESC")
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(simpleJsonResponse("500 Internal Server Error", "Error retrieving the log entries"))
			return
		}

		// FIXME: mostly just returning log entries in raw format for testing and to get the compiler to stop complaining that I'm not using 'logs'
		resBytes, err := json.Marshal(logs)
		if err != nil {
			log.Println("Error marshalling log entries --")
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(simpleJsonResponse("500 Internal Server Error", "Error marshalling log entires"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)

		// pass LogEntry array to file render? gotta get that data into the ui
	} else if r.Method == "POST" {
		// retrieve logs matching search query
		r.ParseForm()
		// sanitize the shit out of that query
		// then it's more or less just like above
		// maybe add some stuff related to search effectiveness, close matches, etc
	}
	// gorilla mux should return a 404 if method doesn't match GET or POST
}

func getLogByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // This route won't be run if vars["id"] == ""
	logs, err := logQuery(fmt.Sprintf(`SELECT * FROM logs WHERE id = %s ORDER BY timestamp DESC`, vars["id"]))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(simpleJsonResponse("500 Internal Server Error", "Error retrieving the log entries"))
		return
	}

	// FIXME: mostly just returning log entries in raw format for testing and to get the compiler to stop complaining that I'm not using 'logs'
	resBytes, err := json.Marshal(logs[0]) // Should only be one log matching id query
	if err != nil {
		log.Println("Error marshalling log entries --")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(simpleJsonResponse("500 Internal Server Error", "Error marshalling log entires"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)
}

func corsProxyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: serve resource at r.URL.String()
	// TODO: gonna need some serious error handling
}

// Retrieve logs from colloiddb
func logQuery(query string) ([]LogEntry, error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error querying database --")
		return nil, err
	}

	var ret []LogEntry
	for rows.Next() {
		var (
			id      int
			timeStr string
			title   string
			entry   string
		)
		err = rows.Scan(&id, &timeStr, &title, &entry)
		if err != nil {
			log.Println("Error scanning DB row --")
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Println("Error parsing time string from DB --")
			return nil, err
		}

		ret = append(ret, LogEntry{Id: id, Timestamp: parsedTime, Title: title, Entry: entry})
	}
	return ret, nil
}

// Overkill in most cases
func errPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: organize pubip stuff into util package
func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)

	// Initialize DB
	database, err := sql.Open("postgres", dbAddress)
	if err != nil {
		log.Println("Error initializing database --")
		panic(err)
	}
	db = database

	// Print public ip address and ListenAndServe port
	cmd := exec.Command("wget", `http://ipinfo.io/ip`, "-qO", "-")
	std_stream, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Error initializing pipe to stdout --")
		log.Println(err)
		return // Not a big deal if public ip isn't printed
	}
	if err := cmd.Start(); err != nil {
		log.Println("Error executing bash command --")
		log.Println(err)
		return // Same reason as return statement above
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(std_stream)
	pub_ip := strings.TrimSpace(buf.String())
	if pub_ip == "" {
		log.Fatal("Retrived public ip is empty - Are you connected to the internet?")
	}
	log.Println("ListenAndServe at", pub_ip+servPort)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/add", addLogHandler)
	r.HandleFunc("/get", getLogHandler).Methods("GET", "POST")
	r.HandleFunc("/get/{id}", getLogByIdHandler).Methods("GET")
	r.HandleFunc("/cors", corsProxyHandler)

	// FIXME: serve static files using gorilla mux
	http.Handle("/assets", http.FileServer(http.Dir("assets")))

	log.Fatal(http.ListenAndServe(servPort, r))
}
