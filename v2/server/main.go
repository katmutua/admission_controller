package main

import(
  "fmt"
  "net/http"
  "log"
  "html"
)

func processRequests(w http.ResponseWriter, r *http.Request){
    curr_url := html.EscapeString(r.URL.Path)
    fmt.Println(curr_url)

    switch curr_url {
    case "/":
      fmt.Fprintf(w, "Welcome to the HomePage!")
      fmt.Println("Endpoint Hit: homePage")
    case "/validate":
      fmt.Fprintf(w, "Validate endpoint for your kubernetes validation webhook!")
      fmt.Println("Endpoint Hit: Validate")
    default:
      fmt.Fprintf(w, "Sorry, no junk in the server :| !")
    }
}

func handleRequests() {
    http.HandleFunc("/", processRequests)
    http.HandleFunc("/validate", processRequests)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
    handleRequests()
}
