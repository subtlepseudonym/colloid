/* Colloid Server
 * subtlepseudonym
 */

package main

import (
	"net/http"
	"log"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

var (
	las_port string = ":8080"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)

	// Print public ip address and ListenAndServe port
	cmd := exec.Command("wget", `http://ipinfo.io/ip`, "-qO", "-")
	std_stream, err := cmd.StdoutPipe()
	if err != nil { log.Fatal(err) }
	if err := cmd.Start(); err != nil { 
		log.Println("Error retrieving public ip - you may not be connected to the internet")
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(std_stream)
	pub_ip := strings.TrimSpace(buf.String())
	log.Println("ListenAndServe at", pub_ip + las_port)
}

func main() {
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(las_port, nil))
}