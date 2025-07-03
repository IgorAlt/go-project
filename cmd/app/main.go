package main

import (
	"fmt"
	"log"
	"net/http"
	"unrealProject/db"
	"unrealProject/handlers"
	handlers2 "unrealProject/internal/handlers"
	"unrealProject/internal/repository"
	"unrealProject/internal/service"
	"unrealProject/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.Init()

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handlers2.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware)

	r.Get("/", handlers.HomeHandler)
	r.Get("/hello", handlers.HelloHandler)
	r.Route("/users", func(r chi.Router) {
		r.Post("/create", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUserById)
		r.Delete("/{id}", handlers.DeleteUser)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
