package main

import (
	"log"
	"net/http"
	"os"

	"github.com/relief-atlas/backend/db"
	"github.com/relief-atlas/backend/handlers"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()
	log.Println("connected to Oracle")

	h := handlers.New(database)
	mux := http.NewServeMux()

	// ── Disasters ────────────────────────────────────────────
	mux.HandleFunc("GET /api/disasters", h.ListDisasters)
	mux.HandleFunc("POST /api/disasters", h.CreateDisaster)
	mux.HandleFunc("GET /api/disasters/{id}", h.GetDisaster)
	mux.HandleFunc("PUT /api/disasters/{id}", h.UpdateDisaster)
	mux.HandleFunc("DELETE /api/disasters/{id}", h.DeleteDisaster)

	// ── Affected Areas ───────────────────────────────────────
	mux.HandleFunc("GET /api/areas", h.ListAreas)
	mux.HandleFunc("POST /api/areas", h.CreateArea)
	mux.HandleFunc("GET /api/areas/{id}", h.GetArea)
	mux.HandleFunc("PUT /api/areas/{id}", h.UpdateArea)
	mux.HandleFunc("DELETE /api/areas/{id}", h.DeleteArea)

	// ── Shelters ─────────────────────────────────────────────
	mux.HandleFunc("GET /api/shelters", h.ListShelters)
	mux.HandleFunc("POST /api/shelters", h.CreateShelter)
	mux.HandleFunc("GET /api/shelters/{id}", h.GetShelter)
	mux.HandleFunc("PUT /api/shelters/{id}", h.UpdateShelter)
	mux.HandleFunc("DELETE /api/shelters/{id}", h.DeleteShelter)

	// ── Admins ───────────────────────────────────────────────
	mux.HandleFunc("GET /api/admins", h.ListAdmins)
	mux.HandleFunc("POST /api/admins", h.CreateAdmin)
	mux.HandleFunc("POST /api/admins/login", h.AdminLogin)
	mux.HandleFunc("GET /api/admins/{id}", h.GetAdmin)
	mux.HandleFunc("PUT /api/admins/{id}", h.UpdateAdmin)
	mux.HandleFunc("DELETE /api/admins/{id}", h.DeleteAdmin)

	// ── Distributions ────────────────────────────────────────
	mux.HandleFunc("GET /api/distributions", h.ListDistributions)
	mux.HandleFunc("POST /api/distributions", h.CreateDistribution)
	mux.HandleFunc("GET /api/distributions/{id}", h.GetDistribution)
	mux.HandleFunc("PUT /api/distributions/{id}", h.UpdateDistribution)
	mux.HandleFunc("DELETE /api/distributions/{id}", h.DeleteDistribution)

	// ── Users ────────────────────────────────────────────────
	mux.HandleFunc("GET /api/users", h.ListUsers)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("POST /api/users/login", h.UserLogin)
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)
	mux.HandleFunc("PUT /api/users/{id}", h.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", h.DeleteUser)

	// ── Donations ────────────────────────────────────────────
	mux.HandleFunc("GET /api/donations", h.ListDonations)
	mux.HandleFunc("POST /api/donations", h.CreateDonation)
	mux.HandleFunc("GET /api/donations/{id}", h.GetDonation)
	mux.HandleFunc("DELETE /api/donations/{id}", h.DeleteDonation)

	port := getenv("PORT", "8080")
	log.Printf("API listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
