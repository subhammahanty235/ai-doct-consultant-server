package service

import (
	"context"
	"time"

	"github.com/subhammahanty235/medai/internal/db"
	"github.com/subhammahanty235/medai/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentService struct {
	db *db.Database
}

func NewAppointmentService(database *db.Database) *AppointmentService {
	return &AppointmentService{
		db: database,
	}
}

func (s *AppointmentService) BookAppointment(userID primitive.ObjectID, req models.AppointmentRequest, chatSessionID primitive.ObjectID) (*models.Appointment, error) {
	collection := s.db.GetCollection("appointments")

	realDoctorID, err := primitive.ObjectIDFromHex(req.RealDoctorID)
	if err != nil {
		return nil, err
	}

	appointment := models.Appointment{
		UserID:           userID,
		RealDoctorID:     realDoctorID,
		ChatSessionID:    chatSessionID,
		AppointmentDate:  req.AppointmentDate,
		Status:           "pending",
		Symptoms:         req.Symptoms,
		AIRecommendation: req.AIRecommendation,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), appointment)
	if err != nil {
		return nil, err
	}

	appointment.ID = result.InsertedID.(primitive.ObjectID)
	return &appointment, nil
}

func (s *AppointmentService) GetUserAppointments(userID primitive.ObjectID) ([]models.Appointment, error) {
	collection := s.db.GetCollection("appointments")

	cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var appointments []models.Appointment
	if err = cursor.All(context.Background(), &appointments); err != nil {
		return nil, err
	}

	return appointments, nil
}

func (s *AppointmentService) UpdateAppointmentStatus(appointmentID primitive.ObjectID, status string) error {
	collection := s.db.GetCollection("appointments")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": appointmentID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)

	return err
}

func (s *AppointmentService) GetAppointmentByID(appointmentID primitive.ObjectID) (*models.Appointment, error) {
	collection := s.db.GetCollection("appointments")

	var appointment models.Appointment
	err := collection.FindOne(context.Background(), bson.M{"_id": appointmentID}).Decode(&appointment)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}
