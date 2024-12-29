package main

import ( 
	
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	//mux.HandleFunc("/", rootHandler)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}






