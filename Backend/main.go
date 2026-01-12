package main

import (
	"backend-test-mekari/config"
	_ "backend-test-mekari/docs"
	"backend-test-mekari/internal/controllers"
	"backend-test-mekari/internal/repositories"
	"backend-test-mekari/internal/routers"
	"backend-test-mekari/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Expense Management API
// @version 1.0
// @description API untuk tugas teknis Mekari.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <token>
func main() {
	// 1. Database & Migration
	db := config.InitDB()

	// 2. Init Repository, Service, dan Controller
	userRepo := repositories.NewUserRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)

	authSvc := services.NewAuthService(userRepo)
	expenseSvc := services.NewExpenseService(expenseRepo)

	authCtrl := controllers.NewAuthController(authSvc)
	expenseCtrl := controllers.NewExpenseController(expenseSvc)

	// 3. Setup Router
	router := routers.SetupRouter(expenseCtrl, authCtrl)

	// 4. Swagger ke Router (PENTING: Harus masuk ke router utama)
	router.Handle("/swagger/", httpSwagger.WrapHandler)

	// 5. Setup Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 6. Cetak Banner dan Link Aktif
	printBanner(port)

	// 7. Start Server
	log.Fatal(http.ListenAndServe(":"+port, CORSMiddleware(router)))
}

func printBanner(port string) {
	banner := `
    ==================================================
       EXPENSE MANAGEMENT SYSTEM - MEKARI CHALLENGE   
    ==================================================
    `
	fmt.Println(banner)
	fmt.Printf("➜  Local:   http://localhost:%s/\n", port)
	fmt.Printf("➜  Health:  http://localhost:%s/api/health\n", port)
	fmt.Printf("➜  Swagger: http://localhost:%s/swagger/index.html\n", port)
	fmt.Println("==================================================")
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID, X-User-Role")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
