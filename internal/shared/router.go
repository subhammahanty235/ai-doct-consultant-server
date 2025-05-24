package shared

import (
	"github.com/gin-gonic/gin"
	"github.com/subhammahanty235/medai/internal/config"
	"github.com/subhammahanty235/medai/internal/db"
	"github.com/subhammahanty235/medai/internal/handlers"

	"github.com/subhammahanty235/medai/internal/middleware"

	// "github.com/subhammahanty235/medai/internal/handlers"
	"github.com/subhammahanty235/medai/internal/service"
	"github.com/subhammahanty235/medai/internal/utils"
)

func SetupRouter(database *db.Database, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Initialize services
	authService := service.NewAuthService(database, cfg.JWTSecret)
	doctorService := service.NewDoctorService(database)
	appointmentService := service.NewAppointmentService(database)

	// Initialize Gemini client
	geminiClient, err := utils.NewGeminiClient(cfg.GeminiAPIKey)
	if err != nil {
		panic("Failed to initialize Gemini client: " + err.Error())
	}

	chatService := service.NewChatService(database, geminiClient, doctorService)

	// Initialize S3 client
	s3Client, err := utils.NewS3Client(cfg.AWSRegion, cfg.AWSAccessKey, cfg.AWSSecretKey, cfg.S3Bucket)
	if err != nil {
		panic("Failed to initialize S3 client: " + err.Error())
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	doctorHandler := handlers.NewDoctorHandler(doctorService)
	chatHandler := handlers.NewChatHandler(chatService, s3Client)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentService)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/doctors", doctorHandler.GetAllDoctors)
		public.GET("/doctors/real", doctorHandler.GetRealDoctors)
		public.GET("/doctors/:id", doctorHandler.GetDoctorByID)
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// Auth routes
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Chat routes
		protected.POST("/chat/start/:doctorId", chatHandler.StartChat)
		protected.POST("/chat/:sessionId/message", chatHandler.SendMessage)
		protected.POST("/chat/:sessionId/upload", chatHandler.UploadImage)
		protected.GET("/chat/history", chatHandler.GetChatHistory)
		protected.GET("/chat/:sessionId", chatHandler.GetChatSession)

		// Appointment routes
		protected.POST("/appointments", appointmentHandler.BookAppointment)
		protected.GET("/appointments", appointmentHandler.GetUserAppointments)
		protected.PUT("/appointments/:id/status", appointmentHandler.UpdateAppointmentStatus)
		protected.GET("/appointments/:id", appointmentHandler.GetAppointment)
	}

	return r
}
