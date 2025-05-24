package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/subhammahanty235/medai/internal/middleware"
	"github.com/subhammahanty235/medai/internal/models"
	"github.com/subhammahanty235/medai/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}
func (h *AppointmentHandler) BookAppointment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get chat session ID from query parameter (optional)
	chatSessionIDStr := c.Query("chat_session_id")
	var chatSessionID primitive.ObjectID
	if chatSessionIDStr != "" {
		var err error
		chatSessionID, err = primitive.ObjectIDFromHex(chatSessionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat session ID"})
			return
		}
	}

	appointment, err := h.appointmentService.BookAppointment(userID, req, chatSessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}

func (h *AppointmentHandler) GetUserAppointments(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	appointments, err := h.appointmentService.GetUserAppointments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

func (h *AppointmentHandler) UpdateAppointmentStatus(c *gin.Context) {
	appointmentIDStr := c.Param("id")
	appointmentID, err := primitive.ObjectIDFromHex(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.appointmentService.UpdateAppointmentStatus(appointmentID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment status updated successfully"})
}

func (h *AppointmentHandler) GetAppointment(c *gin.Context) {
	appointmentIDStr := c.Param("id")
	appointmentID, err := primitive.ObjectIDFromHex(appointmentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := h.appointmentService.GetAppointmentByID(appointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, appointment)
}
