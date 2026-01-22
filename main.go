package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//go:embed swagger.yaml
var swaggerYAML []byte

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	mux.HandleFunc("/api/produk/", handleProdukDetail)

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	mux.HandleFunc("/api/produk", handleProdukList)

	// GET localhost:8080/api/categories/{id}
	// PUT localhost:8080/api/categories/{id}
	// DELETE localhost:8080/api/categories/{id}
	mux.HandleFunc("/api/categories/", handleCategoryDetail)

	// GET localhost:8080/api/categories
	// POST localhost:8080/api/categories
	mux.HandleFunc("/api/categories", handleCategoryList)

	// localhost:8080/health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// Swagger UI handling
	// 1. Serve swagger.yaml form embedded binary
	mux.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml")
		w.Write(swaggerYAML)
	})

	// 2. Serve Swagger UI HTML
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Kasir API Documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/swagger.yaml',
      dom_id: '#swagger-ui',
    });
  };
</script>
</body>
</html>`)
	})

	// 3. Redirect / to /swagger
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running di port " + port)

	err := http.ListenAndServe(":"+port, corsMiddleware(mux))
	if err != nil {
		fmt.Printf("gagal running server: %v\n", err)
	}
}
