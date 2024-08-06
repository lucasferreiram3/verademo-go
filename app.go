// Will possibly be the main driver file for the app
package main

import (
	"fmt"
	"net/http"
)

func helloPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", helloPage)
	http.ListenAndServe(" ", nil)
}
