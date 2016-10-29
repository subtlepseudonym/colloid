/* Colloid Server
 * subtlepseudonym
 */

package main

import (
	"http"
	"log"
	"os"
	"fmt"
)

func indexHandler(res http.ResponseWriter, req *http.REquest) {
	fmt.Fprintf(w, "This is the index.")
}

func init() {
	// Enable logging to local file
	log.SetOutput(os.Stdout)
}

func main() {
	http.HandlerFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe("8080", nil))
}