/* Colloid Server
 * subtlepseudonym (subtlepseudonym@gmail.com)
 */

package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type LogEntry struct {
	timestamp time.Time
	title string
	entry string
}

const (
	servPort string = ":8080"
	// FIXME: properly set up static ips
	dbAddress string = "postgres://postgres:postgres@192.168.1.4/colloiddb?sslmode=require"
)

var (
	db *sql.DB
)

// TODO: organize handlers into rest package
// TODO: rewrite handlers as REST endpoints
// TODO: serve files on Angular 2 front end
func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: serve bootstrap locally (rather than cdn in index.html)
	http.ServeFile(w, r, "assets/index.html")
}

func addLogHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: insert log into db
}

func getLogHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: retrieve logs from db
	// TODO: POST or GET?
}

func corsProxyHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: serve resource at r.URL.String()
	// TODO: gonna need some serious error handling
}

// Retrieve logs from colloiddb
func getLogs(query string) []LogEntry {
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error querying database --")
		log.Println(err)
	}

	var ret []LogEntry
	for rows.Next() {
		var (
			timeStr string
			title string
			entry string
		)
		err = rows.Scan(&timeStr, &title, &entry)
		if err != nil {
			log.Println("Error scanning DB row --")
			log.Println(err)
		}

		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			log.Println("Error parsing time string from DB --")
			log.Println(err)
		}

		ret = append(ret, LogEntry{ timestamp: parsedTime, title: title, entry: entry })
	}
	return ret
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
	// FIXME: db calls are test code -- to be removed
	rows := getLogs("SELECT * FROM logs ORDER BY timestamp")
	for _, entry := range rows {
		fmt.Println(strings.Join([]string{entry.timestamp.Format(time.Stamp), entry.title, entry.entry}, "\t|\t"))
	}

	// TODO: start using muxer for endpoints
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addLogHandler)
	http.HandleFunc("/get", getLogHandler)
	http.HandleFunc("/cors", corsProxyHandler)

	http.Handle("/assets", http.FileServer(http.Dir("assets")))

	log.Fatal(http.ListenAndServe(servPort, nil))
}
