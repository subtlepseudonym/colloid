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

var (
	lasPort string = ":8080"
	db *sql.DB
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: serve bootstrap locally (rather than cdn in index.html)
	http.ServeFile(w, r, "assets/index.html")
}

func logViewHandler(w http.ResponseWriter, r *http.Request) {
	// TODO serve logs based upon POST data
	// SQL DB running on Elsweyr - colloiddb with table logs
}

func newLogHandler(w http.ResponseWriter, r *http.Request) {
	// TODO serve new log form
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	// TODO receive html encoded text and store to postgres
}

// Retrieve logs from colloiddb
func getLogs(query string) []LogEntry {
	rows, err := db.Query(query)
	checkErr(err)

	var ret []LogEntry
	for rows.Next() {
		var (
			timeStr string
			title string
			entry string
		)
		err = rows.Scan(&timeStr, &title, &entry)
		checkErr(err)

		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		checkErr(err)
		ret = append(ret, LogEntry{ timestamp: parsedTime, title: title, entry: entry })
	}
	return ret
}

// Might be overkill
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)

	// Initialize DB
	database, err := sql.Open("postgres",
		"postgres://postgres:postgres@192.168.1.3/colloiddb?sslmode=require")
	checkErr(err)
	db = database

	// Print public ip address and ListenAndServe port
	cmd := exec.Command("wget", `http://ipinfo.io/ip`, "-qO", "-")
	std_stream, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(std_stream)
	pub_ip := strings.TrimSpace(buf.String())
	if pub_ip == "" {
		log.Fatal("Error retrieving public ip - you may not be connected to the internet")
	}
	log.Println("ListenAndServe at", pub_ip+lasPort)
}

func main() {
	rows := getLogs("SELECT * FROM logs ORDER BY timestamp")
	for _, entry := range rows {
		fmt.Println(strings.Join([]string{entry.timestamp.Format(time.Stamp), entry.title, entry.entry}, " | "))
	}

	http.HandleFunc("/", indexHandler) // webpage
	http.HandleFunc("/view", logViewHandler) // webpage
	http.HandleFunc("/new", newLogHandler) // webpage
	http.HandleFunc("/log", logHandler) // REST endpoint

	http.Handle("/assets", http.FileServer(http.Dir("assets")))

	log.Fatal(http.ListenAndServe(lasPort, nil))
}
