package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	name = "member-service"
	port = 3003
)

type member struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func main() {
	members := []member{
		{
			ID:          "123",
			DisplayName: "Alex",
		},
		{
			ID:          "456",
			DisplayName: "Chandra",
		},
		{
			ID:          "789",
			DisplayName: "Orlando",
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, name)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	// private
	mux.HandleFunc("GET /api/v1/members", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(members)
	})
	mux.HandleFunc("GET /api/v1/members/{member_id}", func(w http.ResponseWriter, r *http.Request) {
		memID := r.PathValue("member_id")
		for _, mem := range members {
			if memID == mem.ID {
				json.NewEncoder(w).Encode(mem)
				return
			}
		}

		json.NewEncoder(w).Encode(members[0])
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("serve http server on port:", port)
	log.Fatalf("cannot start http server: %v\n", srv.ListenAndServe())
}
