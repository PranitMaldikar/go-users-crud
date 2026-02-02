package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go-users-crud/controllers"
	"go-users-crud/services"
)

func main() {

	dsn := "postgres://app:app123@localhost:5433/usersdb?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	mux := http.NewServeMux()
	userController.RegisterRoutes(mux)

	log.Println("API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
