package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kgermando/dentic-support-api/controlleurs/auth"
)

// SetupAuthRoutes configure toutes les routes d'authentification
func SetupAuthRoutes(api fiber.Router) {
	a := api.Group("/auth")

	// Public authentication routes
	a.Post("/register", auth.Register)
	a.Post("/login", auth.Login)
	a.Post("/forgot-password", auth.Forgot)
	a.Get("/verify-reset-token/:token", auth.VerifyResetToken)
	a.Post("/reset/:token", auth.ResetPassword)

	// Protected authentication routes
	// app.Use(middlewares.IsAuthenticated)
	a.Get("/agent", auth.AuthAgent)
	a.Put("/profil/info", auth.UpdateInfo)
	a.Put("/change-password", auth.ChangePassword)
	a.Post("/logout", auth.Logout)
}
