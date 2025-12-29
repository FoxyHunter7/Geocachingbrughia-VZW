package router

import (
	"net/http"
	"time"

	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/config"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/database"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/handlers"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/middleware"
	"github.com/FoxyHunter7/geocachingbrughia-backend/internal/services/email"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(db *database.DB, cfg *config.Config, emailService *email.Service) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize handlers
	h := handlers.New(db, cfg, emailService)

	// Initialize login rate limiter (5 attempts per 15 minutes per IP)
	loginLimiter := middleware.NewRateLimiter(5, 15*time.Minute)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Public API routes
	r.Route("/api", func(r chi.Router) {
		// Rate limiting for public endpoints
		r.Use(chiMiddleware.Throttle(100))

		// Public endpoints (with ETag caching - always revalidate but use cache if unchanged)
		r.With(middleware.CacheControl()).Get("/languages", h.GetLanguages)
		r.With(middleware.CacheControl()).Get("/static", h.GetStaticContent)
		r.With(middleware.CacheControl()).Get("/socials", h.GetSocials)
		r.With(middleware.CacheControl()).Get("/events", h.GetPublicEvents)
		r.With(middleware.CacheControl()).Get("/events/{uuid}", h.GetEventByUUID)
		r.Get("/events/{uuid}/qr-codes", h.GetEventQRCodes)
		r.With(middleware.CacheControl()).Get("/home_events", h.GetHomeEvents)
		r.With(middleware.CacheControl()).Get("/messages", h.GetPublicMessages)
		r.With(middleware.CacheControl()).Get("/geocaches", h.GetPublicGeocaches)

		// Serve uploaded images (public, cached)
		r.Get("/images/*", h.ServeImage)

		// Contact form (no caching, rate limited more strictly)
		r.With(chiMiddleware.Throttle(10)).Post("/contact", h.SubmitContactForm)

		// Authentication (login only - no registration endpoint, admin is seeded on first boot)
		// Rate limited per IP: 5 attempts per 15 minutes
		r.With(middleware.LoginRateLimit(loginLimiter)).Post("/login", h.Login)

		// Protected admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWT.Secret))

			// Auth
			r.Get("/profile", h.GetProfile)
			r.Post("/refresh", h.RefreshToken)
			r.Post("/logout", h.Logout)
			r.Post("/change-password", h.ChangePassword)

			// Image upload
			r.Post("/upload-image", h.UploadImage)

			// Events CRUD
			r.Get("/events", h.GetAdminEvents)
			r.Get("/events/{id}", h.GetEventByID)
			r.Post("/events", h.CreateEvent)
			r.Put("/events/{id}", h.UpdateEvent)
			r.Delete("/events/{id}", h.DeleteEvent)

			// Geocaches CRUD
			r.Get("/geocaches", h.GetAdminGeocaches)
			r.Get("/geocaches/{id}", h.GetGeocacheByID)
			r.Post("/geocaches", h.CreateGeocache)
			r.Put("/geocaches/{id}", h.UpdateGeocache)
			r.Delete("/geocaches/{id}", h.DeleteGeocache)

			// Messages CRUD
			r.Get("/messages", h.GetAdminMessages)
			r.Get("/messages/{id}", h.GetMessageByID)
			r.Post("/messages", h.CreateMessage)
			r.Put("/messages/{id}", h.UpdateMessage)
			r.Delete("/messages/{id}", h.DeleteMessage)

			// Languages CRUD
			r.Get("/languages", h.GetAdminLanguages)
			r.Post("/languages", h.CreateLanguage)
			r.Put("/languages/{code}", h.UpdateLanguage)
			r.Delete("/languages/{code}", h.DeleteLanguage)

			// Static content / translations CRUD
			r.Get("/static", h.GetAdminStaticContent)
			r.Post("/static", h.CreateStaticContent)
			r.Put("/static/{property}", h.UpdateStaticContent)
			r.Delete("/static/{property}", h.DeleteStaticContent)

			// Socials CRUD
			r.Get("/socials", h.GetAdminSocials)
			r.Post("/socials", h.CreateSocial)
			r.Put("/socials/{id}", h.UpdateSocial)
			r.Delete("/socials/{id}", h.DeleteSocial)

			// Contact submissions management
			r.Get("/contacts", h.GetContactSubmissions)
			r.Get("/contacts/{id}", h.GetContactSubmissionByID)
			r.Put("/contacts/{id}", h.UpdateContactSubmission)
			r.Post("/contacts/{id}/notes", h.AddContactNote)
			r.Delete("/contacts/{id}", h.DeleteContactSubmission)

			// User management
			r.Get("/users", h.GetUsers)
			r.Post("/users", h.CreateUser)
			r.Delete("/users/{id}", h.DeleteUser)
			r.Post("/users/{id}/resend-invitation", h.ResendInvitation)
		})
	})

	return r
}
