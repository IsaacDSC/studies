package main

import (
	"log"
	"net/http"
	"os"

	"rulesengine/internal/domain/blocklist"
	"rulesengine/internal/domain/campaigns"
	"rulesengine/internal/domain/operation"
	"rulesengine/internal/middleware"
)

func main() {
	// Initialize repositories
	blocklistRepo := blocklist.NewFileRepository(getEnv("BLOCKLIST_PATH", "./blocklist"))
	campaignsRepo := campaigns.NewFileRepository(getEnv("CAMPAIGNS_PATH", "./campaigns"))

	// Initialize services
	blocklistSvc := blocklist.NewService(blocklistRepo)
	campaignsSvc := campaigns.NewService(campaignsRepo)

	// Initialize engine
	operationEngine := operation.NewEngine(campaignsRepo, blocklistSvc)

	// Initialize handlers
	blocklistHandler := blocklist.NewHandler(blocklistSvc)
	campaignsHandler := campaigns.NewHandler(campaignsSvc)
	operationHandler := operation.NewHandler(operationEngine)

	// Setup routes
	mux := http.NewServeMux()

	// Blocklist routes
	mux.HandleFunc("POST /backoffice/block-list", blocklistHandler.Block)
	mux.HandleFunc("DELETE /backoffice/block-list/{user_id}", blocklistHandler.Unblock)

	// Campaigns routes
	mux.HandleFunc("PATCH /backoffice/campaigns/{campaignName}", campaignsHandler.Upsert)
	mux.HandleFunc("GET /backoffice/campaigns/list", campaignsHandler.List)
	mux.HandleFunc("DELETE /backoffice/campaigns/{campaignName}/{operation}", campaignsHandler.Delete)

	// Operation routes
	mux.HandleFunc("POST /rule-engine/rule/{operation}", operationHandler.Evaluate)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Apply middleware
	handler := middleware.Recovery(middleware.Logging(mux))

	// Start server
	port := getEnv("HTTP_PORT", "8080")
	log.Printf("Server starting on :%s", port)
	log.Printf("Blocklist path: %s", getEnv("BLOCKLIST_PATH", "./blocklist"))
	log.Printf("Campaigns path: %s", getEnv("CAMPAIGNS_PATH", "./campaigns"))

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
