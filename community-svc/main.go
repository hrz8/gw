package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	name = "community-service"
	port = 3002
)

type community struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount uint16 `json:"member_count"`
}

func main() {
	communities := []community{
		{
			ID:          "abc",
			Name:        "Hiking",
			MemberCount: 10,
		},
		{
			ID:          "def",
			Name:        "Cycling",
			MemberCount: 7,
		},
		{
			ID:          "ghi",
			Name:        "Swimming",
			MemberCount: 13,
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, name)
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	// public
	mux.HandleFunc("GET /api/v1/communities", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(communities)
	})
	// private
	mux.HandleFunc("GET /api/v1/communities/{community_id}", func(w http.ResponseWriter, r *http.Request) {
		comID := r.PathValue("community_id")
		for _, com := range communities {
			if comID == com.ID {
				json.NewEncoder(w).Encode(com)
				return
			}
		}

		json.NewEncoder(w).Encode(communities[0])
	})
	mux.HandleFunc("GET /api/v1/communities/{community_id}/members", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
		}{
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
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	fmt.Println("serve http server on port:", port)
	log.Fatalf("cannot start http server: %v\n", srv.ListenAndServe())
}
