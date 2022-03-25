package router

import (
	"net/http"
	"time"
	"web_server/middelwares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"web_server/infrastructure"
	"web_server/user"
)

func RunServer() {
	r := chi.NewRouter()

	//middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	userController := user.NewController(infrastructure.Client, infrastructure.UserCollection, infrastructure.DatabaseName)

	r.Group(func(r chi.Router) {
		r.Post("/user/login", userController.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middelwares.Authentication)
		r.Post("/user", userController.CreateUser)
		r.Get("/user/{id}", userController.GetUser)
		r.Get("/user", userController.FilterUser)
		r.Put("/user/{id}", userController.UpdateUser)
		r.Delete("/user", userController.DeleteUser)
	})
	err := http.ListenAndServe(infrastructure.Domain+":"+infrastructure.Port, r)
	if err != nil {
		infrastructure.ErrLog.Fatalln(err)
	}
}
