package routers

import (
	controllers "github.com/Conding-Student/study_template/pkg/controllers"
	"github.com/Conding-Student/study_template/pkg/controllers/healthchecks"
	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

	//database endpoint without token
	// User registration endpoint
	v1Endpoint.Post("/register", controllers.CreateUser)
	v1Endpoint.Post("/login", controllers.LoginUser)

}

func SetupPublicRoutesB(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)

}
