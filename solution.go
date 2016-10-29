/* Colloid Server
 * subtlepseudonym
 */

package main

import (
	"net/http"
	"log"
	"os"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the index.")
}

func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)
}

func main() {
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}