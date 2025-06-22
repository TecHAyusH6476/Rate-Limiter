package main

import (
	"fmt"
	"log"
	"net/http"

	"ratelimit/ratelimit"
)

func main() {
	rl, err := ratelimit.NewRateLimiter("config.yaml")
	fmt.Printf("RateLimiter error: %+v \n", err)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Message sent!")
	})

	// Wrap with middleware
	handler := rl.Middleware(mux)

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
