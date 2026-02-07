package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

//go:embed swagger.yaml
var swaggerYAML []byte

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

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
	// Initialize Viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("Error reading .env file: %v\n", err)
		}
	}

	// Load configuration
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Set default port if not configured
	if config.Port == "" {
		config.Port = "8080"
	}

	// Debug: print config values
	log.Printf("Config - Port: %s\n", config.Port)
	log.Printf("Config - DBConn: %s\n", config.DBConn)

	// Setup database (make it optional for development)
	var db *sql.DB
	var err error
	if config.DBConn != "" {
		db, err = database.InitDB(config.DBConn)
		if err != nil {
			log.Printf("Warning: Failed to initialize database: %v\n", err)
			log.Println("Continuing without database connection...")
		} else {
			defer db.Close()
		}
	} else {
		log.Println("No DB_CONN configured, running without database")
	}

	// Initialize repositories, services, and handlers

	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	mux := http.NewServeMux()

	// POST localhost:8080/api/checkout
	mux.HandleFunc("/api/checkout", transactionHandler.HandleCheckout)

	// GET localhost:8080/api/produk/{id}
	// PUT localhost:8080/api/produk/{id}
	// DELETE localhost:8080/api/produk/{id}
	mux.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// GET localhost:8080/api/produk
	// POST localhost:8080/api/produk
	mux.HandleFunc("/api/produk", productHandler.HandleProducts)

	// GET localhost:8080/api/categories/{id}
	// PUT localhost:8080/api/categories/{id}
	// DELETE localhost:8080/api/categories/{id}
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	// GET localhost:8080/api/categories
	// POST localhost:8080/api/categories
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)

	// GET localhost:8080/api/report/hari-ini
	mux.HandleFunc("/api/report/hari-ini", reportHandler.HandleTodayReport)

	// GET localhost:8080/api/report?start_date=2026-01-01&end_date=2026-02-01
	mux.HandleFunc("/api/report", reportHandler.HandleDateRangeReport)

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
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
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
	  url: '/swagger.yaml?v=' + Date.now(),
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

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	listenErr := http.ListenAndServe(addr, corsMiddleware(mux))
	if listenErr != nil {
		fmt.Printf("gagal running server: %v\n", listenErr)
	}
}
