package handlers

import (
	"net/http"

	"github.com/subhammahanty235/medai/internal/service"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	doctorService *service.DoctorService
}

func NewDoctorHandler(doctorService *service.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		doctorService: doctorService,
	}
}

func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	doctors, err := h.doctorService.GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctors": doctors})
}

func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	doctorID := c.Param("id")

	doctor, err := h.doctorService.GetDoctorByID(doctorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetRealDoctors(c *gin.Context) {
	specialty := c.Query("specialty")

	var doctors []interface{}
	// var err error

	if specialty != "" {
		realDoctors, err := h.doctorService.GetRealDoctorsBySpecialty(specialty)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, doc := range realDoctors {
			doctors = append(doctors, doc)
		}
	} else {
		realDoctors, err := h.doctorService.GetAllRealDoctors()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, doc := range realDoctors {
			doctors = append(doctors, doc)
		}
	}

	c.JSON(http.StatusOK, gin.H{"doctors": doctors})
}
