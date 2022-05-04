package router

import (
	"net/http"
	"time"
	"web_server/middlewares"
	"web_server/scope/api_permission"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"web_server/infrastructure"
	"web_server/scope/api"
	"web_server/scope/role"
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

	userController := user.NewController(infrastructure.Client)
	roleController := role.NewController(infrastructure.Client)
	apiController := api.NewController(infrastructure.Client)
	permissionController := api_permission.NewController(infrastructure.Client)

	r.Group(func(r chi.Router) {
		r.Post("/user/login", userController.Login)
		r.Post("/user/register", userController.Register)
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authentication)
			r.Use(middlewares.Authorization)
			r.Post("/user", userController.CreateUser)
			r.Get("/user/{id}", userController.GetUser)
			r.Get("/user", userController.FilterUser)
			r.Put("/user/{id}", userController.UpdateUser)
			r.Delete("/user", userController.DeleteUser)
		})

		r.Group(func(r chi.Router) {
			//r.Use(middlewares.Authentication)
			r.Post("/role", roleController.CreateRole)
			r.Get("/role/{id}", roleController.GetRoleById)
			r.Get("/role", roleController.FilterRole)
			r.Put("/role/{id}", roleController.UpdateRole)
			r.Delete("/role", roleController.DeleteRole)
		})

		r.Group(func(r chi.Router) {
			//r.Use(middlewares.Authentication)
			r.Post("/api", apiController.CreateApi)
			r.Get("/api/{id}", apiController.GetApiById)
			r.Get("/api", apiController.FilterApi)
			r.Put("/api/{id}", apiController.UpdateApi)
			r.Delete("/api", apiController.DeleteApi)
		})

		r.Group(func(r chi.Router) {
			//r.Use(middlewares.Authentication)
			r.Post("/casbin_rule", permissionController.CreatePermission)
			r.Get("/casbin_rule/{id}", permissionController.GetPermissionById)
			r.Get("/casbin_rule", permissionController.FilterPermission)
			r.Put("/casbin_rule/{id}", permissionController.UpdatePermission)
			r.Delete("/casbin_rule", permissionController.DeletePermission)
		})

	})

	err := http.ListenAndServe(infrastructure.Domain+":"+infrastructure.Port, r)
	if err != nil {
		infrastructure.ErrLog.Fatalln(err)
	}
}
