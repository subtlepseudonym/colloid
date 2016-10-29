/* Colloid Server
 * subtlepseudonym
 */

package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var (
	lasPort string = ":8080"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Index (currently) pulls boostrap from maxcdn rather than local storage
	http.ServeFile(w, r, "web/index.html")
}

func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)

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
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(lasPort, nil))
}
